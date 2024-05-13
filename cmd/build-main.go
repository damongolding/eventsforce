package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	human "github.com/dustin/go-humanize"

	"github.com/charmbracelet/lipgloss"
	"github.com/chromedp/chromedp"
	"github.com/damongolding/eventsforce/internal/utils"
)

func mainBuild(buildMode bool) error {

	var buildOutput []string

	if buildMode {
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

	// Spin up a headless browser for screenshots
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

		if buildMode {
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
				err = cssProcessor(path, buildMode)
				if err != nil {
					return err
				}
			case ".html", ".htm":
				err = htmlProcessor(path, buildMode)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if buildMode {
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

	if buildMode {
		out := lipgloss.JoinVertical(lipgloss.Left, buildOutput...)
		fmt.Println(out)
	}

	return nil
}
