package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

//go:embed new-template-assets
var newTemplateAssets embed.FS

func init() {
	rootCmd.AddCommand(newCmd)
}

// rootCmd represents the base command when called without any subcommands
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new eventsforce template",

	Run: func(cmd *cobra.Command, args []string) {

		var newTemplateName string
		var createNewTemplateConfim bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Template name").
					Description("The name of the new template").
					Value(&newTemplateName),

				huh.NewConfirm().
					Title("Ready?").
					Value(&createNewTemplateConfim),
			),
		)

		fmt.Println()
		if err := form.Run(); err != nil {
			panic(err)
		}

		if createNewTemplateConfim {
			newTemplatePath := filepath.Join(config.SrcDir, newTemplateName)

			if err := createNewTemplate(newTemplatePath); err != nil {
				panic(err)
			}

			fmt.Println(blueBold(newTemplateName), "has been created in", blueBold(newTemplatePath))
		}
	}}

func createNewTemplate(newTemplatePath string) error {

	if _, err := os.Stat(newTemplatePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(newTemplatePath, 0750); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(newTemplatePath, "index.html"))
		if err != nil {
			return err
		}
		defer f.Close()

		indexfile, err := newTemplateAssets.ReadFile("new-template-assets/index.html")
		if err != nil {
			return err
		}

		f.Write(indexfile)

	} else {
		fmt.Println("Hmmmm looks like that template already exists")
	}

	return nil
}
