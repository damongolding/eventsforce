package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/damongolding/eventsforce/internal/utils"
	minify "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

// addCssIncludes processes css @inlude tags, runs recursively
func addCssIncludes(cssContent string) (string, error) {

	// @include
	re := regexp.MustCompile(`@include\s+['"](?P<include>[^'"]*)['"];`)

	matches := re.FindAllStringSubmatch(cssContent, -1)

	for i, match := range matches {
		cssFileContent, err := os.ReadFile(filepath.Join(config.SrcDir, "_includes", match[1]))
		if err != nil {
			return cssContent, err
		}

		cssContent = strings.ReplaceAll(cssContent, match[0], string(cssFileContent))

		if i+1 == len(matches) {
			return addCssIncludes(cssContent)
		}

	}

	return cssContent, nil
}

func cssProcessor(path string, productionMode bool) error {

	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cssContent := string(fileContent)

	cssContent, err = addCssIncludes(cssContent)
	if err != nil {
		return err
	}

	if productionMode {
		// Minifiy CSS
		if config.BuildOptions.MinifyCSS {
			cssContent, err = minifier.String("text/css", cssContent)
			if err != nil {
				return err
			}
		}
	}

	err = os.WriteFile(path, []byte(cssContent), 0777)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := Tailwind(ctx, path, "--minify"); err != nil {
		return err
	}

	if productionMode {
		fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.RemoveDockerPathPrefix(path)))
	}

	return nil

}
