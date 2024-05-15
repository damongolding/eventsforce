/*
Copyright © 2024 NAME HERE damon@damongolding.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/cmd/build"
	"github.com/damongolding/eventsforce/cmd/new"
	"github.com/damongolding/eventsforce/cmd/serve"
	"github.com/damongolding/eventsforce/cmd/version"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
)

var (
	devPort int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eventsforce",
	Short: "eventsforce template builder",
	Long:  `eventsforce template builder`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(showHeader)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().IntVar(&devPort, "port", 3000, "Port to use for dev server")

	// commands
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(new.NewTemplateCmd)
	rootCmd.AddCommand(build.BuildCmd)
	rootCmd.AddCommand(serve.ServeCmd)

}

func showHeader() {

	GraidentColours := []string{
		"#fec5c5",
		"#feafb0",
		"#fe9c9d",
		"#fe898a",
		"#ea7aa3",
		"#d46dc3",
		"#c061e0",
		"#a750ff",
		"#a953fe",
		"#a953fe",
		"#a953fe",
	}

	var builder strings.Builder

	s := "EVENTSFORCE"

	for i, letter := range s {
		l := lipgloss.NewStyle().PaddingRight(1).Bold(true).Foreground(lipgloss.Color(GraidentColours[i])).Render(string(letter))
		builder.WriteString(l)
	}

	out := lipgloss.JoinHorizontal(lipgloss.Left, builder.String())
	fmt.Println(lipgloss.NewStyle().PaddingLeft(1).BorderForeground(lipgloss.Color("#3b414d")).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderBottom(true).Render(out))

}

func print(message ...string) {
	fmt.Println("[ef]", strings.Join(message, " "))
}

func sectionMessage(message ...string) string {
	sectionOutputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderLeft(true).BorderBackground(lipgloss.Color(utils.GraidentColours[0])).
		PaddingLeft(1).
		Render(strings.Join(message, " "))

	return sectionOutputStyle

}
