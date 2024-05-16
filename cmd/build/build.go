package build

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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

		start := time.Now()
		devPort = cmd.Flag("port").Value.String()

		defer func(start time.Time) {
			if err := recover(); err != nil {
				b := utils.Red("Failed in")
				t := utils.BoldRedUnderline(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
				s := utils.Red("seconds")
				l := utils.BoldYellow("😔")

				fmt.Println("\n", err)

				fmt.Println("\n", b, t, s, l)
				os.Exit(1)
			}
		}(start)

		var wg sync.WaitGroup

		wg.Add(1)

		if err := Build(true, &wg); err != nil {
			panic(err)
		}

		wg.Wait()

		b := utils.Green("Built in")
		t := utils.BoldGreenUnderline(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
		s := utils.Green("seconds")
		l := utils.BoldYellow("⚡")

		fmt.Println("\n", b, t, s, l)

	},
}

func Build(buildMode bool, wg *sync.WaitGroup) error {

	stopScreenshotServer := make(chan bool, 1)
	defer close(stopScreenshotServer)
	defer wg.Done()

	if buildMode {
		if err := config.PrintConfig(); err != nil {
			return err
		}
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

func screenshotTemplate(ctx context.Context, file fs.DirEntry) error {

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(fmt.Sprint(`http://localhost:`, devPort, "/", file.Name(), "/"), 90, &buf)); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(config.BuildDir, file.Name(), "screenshot.png"), buf, 0o644); err != nil {
		return err
	}

	fmt.Println(utils.SectionMessage(utils.Green("Created"), filepath.Join(config.BuildDir, file.Name(), "screenshot.png")))

	return nil
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
