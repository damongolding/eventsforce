package version

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
)

var version string
var tailwindVersion string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",

	Run: func(cmd *cobra.Command, args []string) {

		// versionString := lipgloss.NewStyle().Bold(true).Padding(1, 0, 0, 0).Render("Eventsforce", utils.BoldGreen(version))
		// tailwindVersionString := lipgloss.NewStyle().Bold(true).Padding(1, 0, 0, 0).Render("Tailwind", utils.BoldGreen(tailwindVersion))

		t := table.New().
			Border(lipgloss.HiddenBorder()).
			Rows([]string{"Eventsforce", utils.BoldGreen(version)}, []string{"Tailwind", utils.BoldGreen(tailwindVersion)}).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return lipgloss.NewStyle().PaddingRight(1)
				}
				return lipgloss.NewStyle().Padding(0, 1)
			})

		fmt.Println(t.String())

	},
}
