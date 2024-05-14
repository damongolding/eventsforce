package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

// rootCmd represents the base command when called without any subcommands
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build templates for distribution",
	Long:  `Build templates for distribution on the Eventforce platform`,

	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()

		if err := build(true); err != nil {
			panic(err)
		}

		fmt.Println()

		b := green("Built in")
		t := boldGreenUnderline(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
		s := green("seconds")
		l := boldYellow("⚡")

		fmt.Println(b, t, s, l)

	},
}

func build(buildMode bool) error {

	stopScreenshotServer := make(chan bool, 1)
	defer close(stopScreenshotServer)

	// PRE BUILD
	if err := preBuild(buildMode, stopScreenshotServer); err != nil {
		return err
	}
	// END PREBUILD

	// BUILD
	if err := mainBuild(buildMode); err != nil {
		return err
	}
	// END BUILD

	// POST BUILD
	if err := postBuild(buildMode, stopScreenshotServer); err != nil {
		return err
	}
	// END POST BUILD

	return nil
}

func startScreenshotServer(screenshotServerChannel <-chan bool) {

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", devPort),
	}

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	go func() {
		fmt.Println(sectionMessage("Starting screenshot server"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-screenshotServerChannel

	fmt.Println(sectionMessage("Shutting down screenshot server"))
	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func screenshotTemplate(ctx context.Context, file fs.DirEntry) {

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(fmt.Sprint(`http://localhost:`, devPort, "/", file.Name(), "/"), 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config.BuildDir, file.Name(), "screenshot.png"), buf, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sectionMessage(green("Created"), filepath.Join(config.BuildDir, file.Name(), "screenshot.png")))
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
