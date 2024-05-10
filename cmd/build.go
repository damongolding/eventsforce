package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/internal/utils"
	human "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

// rootCmd represents the base command when called without any subcommands
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build templates for distribution",
	Long:  `Build templates for distribution on the Eventforce platform`,

	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()

		if err := build(true); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println()

		b := green("Built in")
		t := boldGreen(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
		s := green("seconds")
		l := boldYellow("⚡")

		print(b, t, s, l)

	},
}

func build(productionMode bool) error {

	if config.BuildOptions.CleanBuildDir {
		err := utils.CleanBuildDir(config.BuildDir)
		if err != nil {
			return err
		}

		if productionMode {
			print(blue("Cleaning build dir"))
		}
	}

	if productionMode {
		fmt.Println()
		b := lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#fff")).Bold(true).Background(lipgloss.Color("#f07e9b")).Render
		u := lipgloss.NewStyle().Padding(0).Foreground(lipgloss.Color("#fff")).Bold(true).Background(lipgloss.Color("#df73b3")).Render
		i := lipgloss.NewStyle().Padding(0).Foreground(lipgloss.Color("#fff")).Bold(true).Background(lipgloss.Color("#d36cc3")).Render
		l := lipgloss.NewStyle().Padding(0).Foreground(lipgloss.Color("#fff")).Bold(true).Background(lipgloss.Color("#c664d5")).Render
		d := lipgloss.NewStyle().PaddingRight(1).Foreground(lipgloss.Color("#fff")).Bold(true).Background(lipgloss.Color("#b95ce8")).Render
		s := fmt.Sprintf("%s%s%s%s%s", b("B"), u("U"), i("I"), l("L"), d("D"))

		bord := lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
			BorderBottom(true).
			Padding(0).
			Render
		fmt.Println(bord(s))
	}

	fileList, err := os.ReadDir(config.SrcDir)
	if err != nil {
		return err
	}

	for _, file := range fileList {

		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), "_") {
			continue
		}

		// Move files
		err := utils.CopyDir(filepath.Join(config.SrcDir, file.Name()), filepath.Join(config.BuildDir, file.Name()))
		if err != nil {
			return err
		}

		// Add fonts
		if config.BuildOptions.AddFonts {
			err = utils.CopyDir(filepath.Join(config.SrcDir, "_assets", "fonts"), filepath.Join(config.BuildDir, file.Name()))
			if err != nil {
				return err
			}
		}

		// Do build things
		filepath.Walk(config.BuildDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			switch filepath.Ext(path) {
			case ".css":
				cssProcessor(path, productionMode)
				if err != nil {
					return err
				}
			case ".html", ".htm":
				htmlProcessor(path, productionMode)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if productionMode {
			fullPath := filepath.Join(config.BuildDir, file.Name())
			zipPath := filepath.Join(config.BuildDir, file.Name()+".zip")

			zipSize, err := utils.ZipDirectory(fullPath, zipPath)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			err = os.RemoveAll(filepath.Join(config.BuildDir, file.Name()))
			if err != nil {
				fmt.Println(err)
				return nil
			}

			print(green("created"), zipPath, "("+human.BigBytes(human.BigByte.SetInt64(zipSize))+")")
		}

	}

	return nil
}
