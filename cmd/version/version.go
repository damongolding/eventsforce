package version

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
)

var version string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",

	Run: func(cmd *cobra.Command, args []string) {
		versionString := lipgloss.NewStyle().Bold(true).Padding(1, 0).Render("You are currently running", utils.BoldGreen(version))
		fmt.Println(versionString)
	},
}
