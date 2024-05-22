package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"

	"github.com/bep/godartsass/v2"
	"github.com/damongolding/eventsforce/internal/utils"
)

type importResolver struct{}

func (i importResolver) CanonicalizeURL(url string) (string, error) {
	return fmt.Sprintf("file://%s", url), nil
}

func (t importResolver) Load(url string) (godartsass.Import, error) {

	includeDir := filepath.Join(config.SrcDir, "_includes")

	filePath := strings.Replace(url, "file://", "", -1)
	filePath = filepath.Join(includeDir, filePath)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return godartsass.Import{}, err
	}

	sourceSyntax := godartsass.SourceSyntaxSCSS
	switch filepath.Ext(url) {
	case ".scss":
		sourceSyntax = godartsass.SourceSyntaxSCSS
	case ".sass":
		sourceSyntax = godartsass.SourceSyntaxSASS
	case ".css":
		sourceSyntax = godartsass.SourceSyntaxCSS
	}

	return godartsass.Import{Content: string(b), SourceSyntax: sourceSyntax}, err
}

func sassProcessor(path string, productionMode bool) error {

	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)

	sassFileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(utils.SectionErrorMessage(err.Error()))
		return err
	}

	args := godartsass.Args{
		Source:         string(sassFileContent),
		ImportResolver: importResolver{},
	}

	var options godartsass.Options

	if utils.RunningInDocker() {
		options.DartSassEmbeddedFilename = "/dart-sass/sass"
	}

	transpiler, err := godartsass.Start(options)
	if err != nil {
		return err
	}

	defer transpiler.Close()

	result, err := transpiler.Execute(args)
	if err != nil {
		return err
	}

	if productionMode {
		// Minifiy CSS
		if config.BuildOptions.MinifyCSS {
			result.CSS, err = minifier.String("text/css", result.CSS)
			if err != nil {
				return err
			}
		}
	}

	err = os.WriteFile(path, []byte(result.CSS), 0666)
	if err != nil {
		return err
	}

	err = os.Rename(path, strings.Replace(path, ".scss", ".css", -1))
	if err != nil {
		return err
	}

	return nil

	// minifier := minify.New()
	// minifier.AddFunc("text/css", css.Minify)

	// fileContent, err := os.ReadFile(path)
	// if err != nil {
	// 	return err
	// }

	// cssContent := string(fileContent)

	// // @include
	// re := regexp.MustCompile(`@include\s+['"](?P<include>[^'"]*)['"];`)

	// matches := re.FindAllStringSubmatch(cssContent, -1)

	// for _, match := range matches {
	// 	cssFileContent, err := os.ReadFile(filepath.Join(config.SrcDir, "_includes", "css", match[1]))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	cssContent = strings.ReplaceAll(cssContent, match[0], string(cssFileContent))

	// }

	// if productionMode {
	// 	// Minifiy CSS
	// 	if config.BuildOptions.MinifyCSS {
	// 		cssContent, err = minifier.String("text/css", cssContent)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// err = os.WriteFile(path, []byte(cssContent), 0666)
	// if err != nil {
	// 	return err
	// }

	return nil
}
