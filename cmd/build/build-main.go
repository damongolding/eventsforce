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
	"golang.org/x/sync/errgroup"
)

func mainBuild(buildMode bool) error {

	defer func() interface{} {
		if err := recover(); err != nil {
			fmt.Println("\n", err)
			return err
		}
		return nil
	}()

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

	buildErrors, ctx := errgroup.WithContext(ctx)

	if err := filepath.Walk(config.BuildDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {

		case ".css":
			buildErrors.Go(func() error {
				err := cssProcessor(path, buildMode)
				return err
			})

		case ".html", ".htm":
			buildErrors.Go(func() error {
				err := htmlProcessor(path, buildMode)
				return err
			})
		}

		return nil
	}); err != nil {
		return err
	}

	if err := buildErrors.Wait(); err != nil {
		return err
	}

	return nil
}
