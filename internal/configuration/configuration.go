package configuration

import (
	"fmt"
	"os"

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
	ZipDirs          bool `mapstructure:"zipDirs"`
}

type Config struct {
	ConfigUsed   string
	SrcDir       string       `mapstructure:"srcDir"`
	BuildDir     string       `mapstructure:"buildDir"`
	BuildOptions BuildOptions `mapstructure:"buildOptions"`
}

func NewConfig() *Config {

	config, err := initConfig()
	if err != nil {
		fmt.Println(utils.SectionErrorMessage(err.Error()))
		os.Exit(1)
	}

	return config
}

// InitConfig reads in config file and ENV variables if set.
func initConfig() (*Config, error) {
	var config Config

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	v := viper.New()

	// Search config in home directory with name ".eventsforce" (without extension).
	v.AddConfigPath(home)
	v.SetConfigFile("config.json")

	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err != nil {
		return &config, err
	}

	if err := v.Unmarshal(&config); err != nil {
		return &config, err

	}

	config.ConfigUsed = v.ConfigFileUsed()

	return &config, nil
}

func (c *Config) PrintConfig() error {
	s, err := utils.OutputSectionStyling("C", "O", "N", "F", "I", "G")
	if err != nil {
		return err
	}
	fmt.Println(s)
	fmt.Println(utils.SectionMessage("Using config file:", utils.BlueBold(c.ConfigUsed)))

	return nil
}
