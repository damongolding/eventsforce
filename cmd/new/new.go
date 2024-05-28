package new

import (
	"embed"
	"errors"
	"fmt"
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
			}

			fmt.Println(utils.BlueBold(newTemplateName), "has been created in", utils.BlueBold(newTemplatePath))
		}
	}}

func createFile(path, fileName, newFileName string) error {
	// Create HTML file
	f, err := os.Create(filepath.Join(path, newFileName))
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

func createNewTemplate(newTemplatePath string, useTailwind bool) error {

	if _, err := os.Stat(newTemplatePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(newTemplatePath, 0755); err != nil {
			return err
		}

		// Create HTML file
		if err := createFile(newTemplatePath, "index.html", "index.html"); err != nil {
			return err
		}

		if useTailwind {
			// Create CSS file
			if err := createFile(newTemplatePath, "style-tailwind.css", "style.css"); err != nil {
				return err
			}
		} else {
			// Create CSS file
			if err := createFile(newTemplatePath, "style.css", "style.css"); err != nil {
				return err
			}
		}

	} else {
		fmt.Println("Hmmmm looks like that template already exists")
	}

	return nil
}
