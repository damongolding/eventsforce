package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/chromedp/chromedp"
	"github.com/damongolding/eventsforce/internal/utils"
	human "github.com/dustin/go-humanize"
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
			fmt.Println(err)
			return
		}

		fmt.Println()

		b := green("Built in")
		t := boldGreen(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
		s := green("seconds")
		l := boldYellow("⚡")

		fmt.Println(b, t, s, l)

	},
}

func build(productionMode bool) error {

	var buildOutput []string

	if config.BuildOptions.CleanBuildDir {
		err := utils.CleanBuildDir(config.BuildDir)
		if err != nil {
			return err
		}

		if productionMode {
			preBuildSectionTitle, err := utils.OutputStyling("P", "R", "E", " ", "B", "U", "I", "L", "D")
			if err != nil {
				return err
			}
			fmt.Println(preBuildSectionTitle)
			fmt.Println(sectionMessage("Cleaning build dir"))

			if err := build(false); err != nil {
				return err
			}

			go startScreenshotServer()
			fmt.Println(sectionMessage("Started screenshot server"))

		}
	}

	if productionMode {
		s, err := utils.OutputStyling("B", "U", "I", "L", "D")
		if err != nil {
			return err
		}
		buildOutput = append(buildOutput, s)

	}

	fileList, err := os.ReadDir(config.SrcDir)
	if err != nil {
		return err
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	defer cancel()

	for _, file := range fileList {

		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), "_") {
			continue
		}

		if productionMode {
			pageScreenshot(ctx, file)
		}

		// Move files
		err := utils.CopyDir(filepath.Join(config.SrcDir, file.Name()), filepath.Join(config.BuildDir, file.Name()))
		if err != nil {
			return err
		}

		// Add fonts
		if config.BuildOptions.AddFonts {
			err = utils.CopyDir(filepath.Join(config.SrcDir, "_assets", "fonts"), filepath.Join(config.BuildDir, file.Name()))
			if err != nil {
				return err
			}
		}

		// Do build things
		filepath.Walk(config.BuildDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			switch filepath.Ext(path) {
			case ".css":
				err = cssProcessor(path, productionMode)
				if err != nil {
					return err
				}
			case ".html", ".htm":
				err = htmlProcessor(path, productionMode)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if productionMode {
			fullPath := filepath.Join(config.BuildDir, file.Name())
			zipPath := filepath.Join(config.BuildDir, file.Name()+".zip")

			zipSize, err := utils.ZipDirectory(fullPath, zipPath)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			err = os.RemoveAll(filepath.Join(config.BuildDir, file.Name()))
			if err != nil {
				fmt.Println(err)
				return nil
			}

			c := fmt.Sprintf("%s %s %s", green("created"), zipPath, "("+human.BigBytes(human.BigByte.SetInt64(zipSize))+")")

			buildOutput = append(buildOutput, sectionMessage(c))

		}

	}

	if productionMode {
		out := lipgloss.JoinVertical(lipgloss.Left, buildOutput...)
		fmt.Println(out)
	}

	return nil
}

func startScreenshotServer() {

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	err := http.ListenAndServe(fmt.Sprintf(":%d", devPort), nil)
	if err != nil {
		panic(err)
	}
}

func pageScreenshot(ctx context.Context, file fs.DirEntry) {

	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(fmt.Sprint(`http://localhost:`, devPort, "/", file.Name(), "/"), 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config.BuildDir, file.Name(), "screenshot.png"), buf, 0o644); err != nil {
		log.Fatal(err)
	}
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
