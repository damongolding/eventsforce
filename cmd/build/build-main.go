package build

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/damongolding/eventsforce/internal/utils"
)

func mainBuild(buildMode bool) error {

	if buildMode {
		s, err := utils.OutputSectionStyling("B", "U", "I", "L", "D")
		if err != nil {
			return err
		}
		fmt.Println(s)
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
			screenshotTemplate(ctx, file)
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

			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Added"), "fonts to", filepath.Join(config.BuildDir, file.Name())))
			}
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
			if err := cssProcessor(path, buildMode); err != nil {
				return err
			}
			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), path))
			}

		case ".html", ".htm":
			if err := htmlProcessor(path, buildMode); err != nil {
				return err
			}
			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), path))
			}

		}

		return nil
	})

	return nil
}
