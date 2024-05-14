package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/damongolding/eventsforce/internal/utils"
	human "github.com/dustin/go-humanize"
)

func postBuild(buildMode bool, stopScreenshotServer chan<- bool) error {

	if buildMode {
		s, err := utils.OutputStyling("P", "O", "S", "T", " ", "B", "U", "I", "L", "D")
		if err != nil {
			return err
		}
		fmt.Println(s)
	}

	stopScreenshotServer <- true

	fileList, err := os.ReadDir(config.BuildDir)
	if err != nil {
		return err
	}

	for _, file := range fileList {

		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), "_") {
			continue
		}

		if buildMode {
			fullPath := filepath.Join(config.BuildDir, file.Name())
			zipPath := filepath.Join(config.BuildDir, file.Name()+".zip")

			zipSize, err := utils.ZipDirectory(fullPath, zipPath)
			if err != nil {
				return err
			}

			err = os.RemoveAll(filepath.Join(config.BuildDir, file.Name()))
			if err != nil {
				return err
			}

			c := fmt.Sprintf("%s %s %s", green("created"), zipPath, "("+human.BigBytes(human.BigByte.SetInt64(zipSize))+")")

			fmt.Println(sectionMessage(c))

		}

	}

	return nil

}
