/*
Copyright © 2024 NAME HERE damon@damongolding.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BuildOptions struct {
	CleanBuildDir    bool `mapstructure:"cleanBuildDir"`
	AddFonts         bool `mapstructure:"addFonts"`
	MinifyHTML       bool `mapstructure:"minifyHTML"`
	KeepCommentsHTML bool `mapstructure:"keepCommentsHTML"`
	MinifyCSS        bool `mapstructure:"minifyCSS"`
}

type Config struct {
	SrcDir       string       `mapstructure:"srcDir"`
	BuildDir     string       `mapstructure:"buildDir"`
	BuildOptions BuildOptions `mapstructure:"buildOptions"`
}

var (
	cfgFile   string
	config    Config
	buildMode bool
	devMode   bool
	devPort   int

	// Colours
	whiteBold          = lipgloss.NewStyle().Bold(true).Render
	green              = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render
	boldGreen          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575")).Render
	boldGreenUnderline = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575")).Underline(true).Render
	yellow             = lipgloss.NewStyle().Foreground(lipgloss.Color("227")).Render
	boldYellow         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("227")).Render
	blue               = lipgloss.NewStyle().Foreground(lipgloss.Color("#9aedff")).Render
	blueBold           = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9aedff")).Render
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

	cobra.OnInitialize(sayHello)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().IntVar(&devPort, "port", 3000, "Port to use for dev server")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sayHello() {

	fmt.Println(os.Getenv("ass"))

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

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".eventsforce" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigFile("config.json")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		s, err := utils.OutputStyling("C", "O", "N", "F", "I", "G")
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
		fmt.Println(sectionMessage("Using config file:", blueBold(viper.ConfigFileUsed())))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic("Missing config file")
	}
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
