package version

import (
	"fmt"
	"os/exec"

	"github.com/bep/godartsass/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
)

var version string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",

	Run: func(cmd *cobra.Command, args []string) {

		versionString := lipgloss.NewStyle().Bold(true).Padding(1, 0, 0, 0).Render("Eventsforce", utils.BoldGreen(version))
		fmt.Println(versionString)

		if utils.RunningInDocker() {
			sassVersion, err := godartsass.Version("/dart-sass/sass")
			if err != nil {
				panic(err)
			}
			sassVersionString := lipgloss.NewStyle().Bold(true).Padding(0).Render("Dart SASS  ", utils.BoldGreen(sassVersion.CompilerVersion))
			fmt.Println(sassVersionString)

		} else {
			sassCmd := exec.Command("sass", "--version")

			sassVersion, err := sassCmd.Output()
			if err != nil {
				panic(err)
			}

			sassVersionString := lipgloss.NewStyle().Bold(true).Padding(0).Render("Dart SASS  ", utils.BoldGreen(string(sassVersion)))
			fmt.Println(sassVersionString)

		}

	},
}
