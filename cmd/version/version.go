package version

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var version string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lipgloss.NewStyle().Bold(true).Padding(1, 0).Render(version))
	},
}
