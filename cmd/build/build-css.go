package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

	start := time.Now()

	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)

	usingTailwaind := false

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cssContent := string(fileContent)

	cssContent, err = addCssIncludes(cssContent)
	if err != nil {
		return err
	}

	if strings.Contains(cssContent, "@tailwind") {
		usingTailwaind = true
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

	if usingTailwaind {
		ctx := context.Background()
		if err := Tailwind(ctx, path, "--minify"); err != nil {
			return err
		}
	}

	if productionMode {

		done := fmt.Sprintf("[%.2fs]", time.Since(start).Seconds())

		if usingTailwaind {
			fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.Blue(done), utils.RemoveDockerPathPrefix(path), "(tailwind)"))
		} else {
			fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.Blue(done), utils.RemoveDockerPathPrefix(path)))
		}
	}

	return nil

}
