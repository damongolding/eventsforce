package new

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/damongolding/eventsforce/internal/configuration"
	"github.com/damongolding/eventsforce/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	//go:embed new-template-assets
	newTemplateAssets embed.FS

	config configuration.Config
)

func init() {
	config = *configuration.NewConfig()
}

type CssTemplate struct {
	Tailwind bool
}

// rootCmd represents the base command when called without any subcommands
var NewTemplateCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new eventsforce template",

	Run: func(cmd *cobra.Command, args []string) {

		var newTemplateName string
		var createNewTemplateConfim bool
		var useTailwind bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Template name").
					Description("The name of the new template").
					Value(&newTemplateName),

				huh.NewConfirm().
					Title("Use Tailwind?").
					Value(&useTailwind),

				huh.NewConfirm().
					Title("Ready?").
					Value(&createNewTemplateConfim),
			),
		)

		fmt.Println()
		if err := form.Run(); err != nil {
			fmt.Println(utils.SectionErrorMessage(err.Error()))
			defer os.Exit(1)
		}

		if createNewTemplateConfim {
			newTemplatePath := filepath.Join(config.SrcDir, newTemplateName)

			if err := createNewTemplate(newTemplatePath, useTailwind); err != nil {
				fmt.Println(utils.SectionErrorMessage(err.Error()))
				defer os.Exit(1)
			} else {
				fmt.Println(utils.BlueBold(newTemplateName), "has been created in", utils.BlueBold(newTemplatePath))
			}

		}
	}}

func createHTML(path, fileName string) error {
	// Create HTML file
	f, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	fileContents, err := newTemplateAssets.ReadFile("new-template-assets/" + fileName)
	if err != nil {
		return err
	}

	f.Write(fileContents)

	return nil
}

func createCSS(path, filename string, useTailwind bool) error {
	f, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	fileContents, err := newTemplateAssets.ReadFile("new-template-assets/style.css")
	if err != nil {
		return err
	}

	tmpl, err := template.New("style").Parse(string(fileContents))
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, CssTemplate{Tailwind: useTailwind})
	if err != nil {
		return err
	}
	return nil
}

func createNewTemplate(newTemplatePath string, useTailwind bool) error {

	if _, err := os.Stat(newTemplatePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(newTemplatePath, 0755); err != nil {
			return err
		}

		// Create HTML file
		if err := createHTML(newTemplatePath, "index.html"); err != nil {
			return err
		}

		// Create CSS
		if err := createCSS(newTemplatePath, "style.css", useTailwind); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Hmmmm looks like that template already exists")
	}

	return nil
}
