package build

import (
	"fmt"

	"github.com/damongolding/eventsforce/internal/utils"
)

// preBuild Things to do beofre the build proccess starts
func preBuild(buildMode bool, stopScreenshotServer <-chan bool) error {

	// Show section title on build
	if buildMode {
		preBuildSectionTitle, err := utils.OutputSectionStyling("P", "R", "E", " ", "B", "U", "I", "L", "D")
		if err != nil {
			return err
		}
		fmt.Println(preBuildSectionTitle)

		go startScreenshotServer(stopScreenshotServer)

	}

	if config.BuildOptions.CleanBuildDir {
		err := utils.CleanBuildDir(config.BuildDir)
		if err != nil {
			return err
		}

		if buildMode {
			fmt.Println(utils.SectionMessage("Cleaning build dir"))
		}
	}

	if buildMode {
		if err := Build(false); err != nil {
			return err
		}

	}

	return nil
}
