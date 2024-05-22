package configuration

import (
	"fmt"
	"io/fs"
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

// findConfigFile tries to find the ocnfig file. Looking in the obvious places first
func findConfigFile(pwd string) (string, error) {

	var foundPath string

	if _, err := os.Stat("/templates/config.json"); err == nil {
		return "/templates/config.json", nil
	} else if _, err := os.Stat("./config.json"); err == nil {
		return "./config.json", nil
	}

	// Lets GO fishing...
	err := filepath.Walk(pwd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "config.json" {
			foundPath = path
			return nil
		}

		return nil
	})

	err = filepath.Walk(filepath.Dir(pwd), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "config.json" {
			foundPath = path
			return nil
		}

		return nil
	})

	return foundPath, err
}

// InitConfig reads in config file and ENV variables if set.
func initConfig() (*Config, error) {
	var config Config

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	v := viper.New()

	v.AddConfigPath(home)

	exe, err := os.Executable()
	if err != nil {
		return &config, err
	}

	pwd := filepath.Dir(exe)
	pwd, err = filepath.Abs(pwd)
	if err != nil {
		return &config, err
	}

	configPath, err := findConfigFile(pwd)
	if err != nil || configPath == "" {
		return &config, err
	}

	v.SetConfigFile(configPath)

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
	} else {
		system := fmt.Sprintf("(%s/%s)", runtime.GOOS, runtime.GOARCH)
		fmt.Println(utils.SectionMessage("Running", utils.BlueBold("natively"), utils.BlueBold(system)))
	}

	fmt.Println(utils.SectionMessage("Using config file:", utils.BlueBold(c.ConfigUsed)))

	return nil
}
