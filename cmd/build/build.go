package build

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
	"github.com/damongolding/eventsforce/internal/configuration"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
)

var (
	config  configuration.Config
	devPort string
)

func init() {
	config = *configuration.NewConfig()
}

// rootCmd represents the base command when called without any subcommands
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build templates for distribution",
	Long:  `Build templates for distribution on the Eventforce platform`,

	Run: func(cmd *cobra.Command, args []string) {

		devPort = cmd.Flag("port").Value.String()

		start := time.Now()

		if err := Build(true); err != nil {
			panic(err)
		}

		fmt.Println()

		b := utils.Green("Built in")
		t := utils.BoldGreenUnderline(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
		s := utils.Green("seconds")
		l := utils.BoldYellow("⚡")

		fmt.Println(b, t, s, l)

	},
}

func Build(buildMode bool) error {

	stopScreenshotServer := make(chan bool, 1)
	defer close(stopScreenshotServer)

	if buildMode {
		config.PrintConfig()
	}

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
		Addr: fmt.Sprintf(":%s", devPort),
	}

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	go func() {
		fmt.Println(utils.SectionMessage("Starting screenshot server"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-screenshotServerChannel

	fmt.Println(utils.SectionMessage("Shutting down screenshot server"))
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

	fmt.Println(utils.SectionMessage(utils.Green("Created"), filepath.Join(config.BuildDir, file.Name(), "screenshot.png")))
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
