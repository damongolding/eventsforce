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

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
	)
	defer cancel()

	for _, file := range fileList {

		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), "_") {
			continue
		}

		fullFilePath := filepath.Join(config.BuildDir, file.Name())

		if buildMode {
			err = screenshotTemplate(ctx, file)
			if err != nil {
				return err
			}
		}

		// Move files
		err := utils.CopyDir(filepath.Join(config.SrcDir, file.Name()), fullFilePath)
		if err != nil {
			return err
		}

		// Add fonts
		if config.BuildOptions.AddFonts {
			err = utils.CopyDir(filepath.Join(config.SrcDir, "_assets", "fonts"), fullFilePath)
			if err != nil {
				return err
			}

			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Added"), "fonts to", utils.RemoveDockerPathPrefix(fullFilePath)))
			}
		}

	}
	// Do build things
	err = filepath.Walk(config.BuildDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".scss":
			if err := sassProcessor(path, buildMode); err != nil {
				return err
			}
			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.RemoveDockerPathPrefix(path)))
			}
		case ".css":
			if err := cssProcessor(path, buildMode); err != nil {
				return err
			}
			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.RemoveDockerPathPrefix(path)))
			}

		case ".html", ".htm":
			if err := htmlProcessor(path, buildMode); err != nil {
				return err
			}
			if buildMode {
				fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.RemoveDockerPathPrefix(path)))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
