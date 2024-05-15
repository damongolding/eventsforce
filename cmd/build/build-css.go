package build

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

func cssProcessor(path string, productionMode bool) error {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cssContent := string(fileContent)

	// @include
	re := regexp.MustCompile(`@include\s+['"](?P<include>[^'"]*)['"];`)

	matches := re.FindAllStringSubmatch(cssContent, -1)

	for _, match := range matches {
		cssFileContent, err := os.ReadFile(filepath.Join(config.SrcDir, "_includes", "css", match[1]))
		if err != nil {
			return err
		}

		cssContent = strings.ReplaceAll(cssContent, match[0], string(cssFileContent))

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

	err = os.WriteFile(path, []byte(cssContent), 0666)
	if err != nil {
		return err
	}

	return nil
}
