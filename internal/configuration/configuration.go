package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

	v.AddConfigPath(home)

	if utils.RunningInDocker() {
		v.SetConfigFile("/eventsforce/config.json")
	} else {
		v.SetConfigFile("config.json")
	}

	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err != nil {
		return &config, err
	}

	if err := v.Unmarshal(&config); err != nil {
		return &config, err

	}

	config.ConfigUsed = v.ConfigFileUsed()

	if utils.RunningInDocker() {
		config.SrcDir = filepath.Join("/templates", config.SrcDir)
		config.BuildDir = filepath.Join("/templates", config.BuildDir)
	}

	return &config, nil
}

func (c *Config) PrintConfig() error {
	s, err := utils.OutputSectionStyling("C", "O", "N", "F", "I", "G")
	if err != nil {
		return err
	}
	fmt.Println(s)

	if utils.RunningInDocker() {
		system := fmt.Sprintf("(%s/%s)", runtime.GOOS, runtime.GOARCH)
		fmt.Println(utils.SectionMessage("Running via", utils.BlueBold("🐳 Docker"), utils.BlueBold(system)))
	}

	fmt.Println(utils.SectionMessage("Using config file:", utils.BlueBold(c.ConfigUsed)))

	return nil
}
