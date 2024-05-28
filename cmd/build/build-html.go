package build

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/damongolding/eventsforce/internal/utils"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func htmlProcessor(path string, productionMode bool) error {

	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
		KeepQuotes:       true,
		KeepComments:     config.BuildOptions.KeepCommentsHTML,
	})

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	htmlContent := string(fileContent)

	if !productionMode {
		//  [[ef_tags]]
		re := regexp.MustCompile(`\[\[(?P<ef_tag>.*)\]\]`)
		matches := re.FindAllStringSubmatch(string(htmlContent), -1)

		for _, match := range matches {
			// [[PATH]] tag specifically
			if strings.TrimSpace(match[1]) == "PATH" {
				htmlContent = strings.Replace(htmlContent, match[0], "./", -1)
				continue
			}

			htmlFileContent, err := os.ReadFile(filepath.Join(config.SrcDir, "_includes", "html", strings.TrimSpace(match[1])+".html"))
			if err != nil {
				fmt.Println(utils.SectionErrorMessage(err.Error()))
				continue
			}

			htmlContent = strings.ReplaceAll(htmlContent, match[0], string(htmlFileContent))

		}
		// Add live reload tag
		htmlContent = strings.Replace(htmlContent, "</body>", "<script src=\"http://localhost:35729/livereload.js\"></script></body>", -1)
	}

	if productionMode {
		htmlContent = strings.Replace(string(fileContent), "<script src=\"http://localhost:35729/livereload.js\"></script>", "", -1)

		// Minify HTML
		if config.BuildOptions.MinifyHTML {
			htmlContent, err = minifier.String("text/html", htmlContent)
			if err != nil {
				return err
			}
		}
	}

	err = os.WriteFile(path, []byte(htmlContent), 0666)
	if err != nil {
		return err
	}

	if productionMode {
		fmt.Println(utils.SectionMessage(utils.Green("Proccessed"), utils.RemoveDockerPathPrefix(path)))
	}

	return nil

}
