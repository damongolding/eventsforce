package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Graident Colurs
	GraidentColours = []string{
		"#f07e9b",
		"#df73b3",
		"#d36cc3",
		"#c664d5",
		"#b95ce8",
		"#b55aef",
		"#a953fe",
		"#a953fe",
		"#a953fe",
		"#a953fe",
		"#a953fe",
		"#a953fe",
		"#a953fe",
	}

	// Colours
	WhiteBold          = lipgloss.NewStyle().Bold(true).Render
	Green              = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render
	BoldGreen          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575")).Render
	BoldGreenUnderline = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575")).Underline(true).Render
	Yellow             = lipgloss.NewStyle().Foreground(lipgloss.Color("227")).Render
	BoldYellow         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("227")).Render
	Blue               = lipgloss.NewStyle().Foreground(lipgloss.Color("#9aedff")).Render
	BlueBold           = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9aedff")).Render
	Red                = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5658")).Render
	BoldRedUnderline   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5658")).Underline(true).Render
)

func RunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func OutputSectionStyling(in ...string) (string, error) {

	var builder strings.Builder

	lettBaseStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fff")).Bold(true).Padding(0)

	if len(in) == 0 || len(in) > len(GraidentColours) {
		return "", fmt.Errorf("Too many letters")
	}

	for i, s := range in {
		var letterStyle func(...string) string

		if i == 0 {
			letterStyle = lettBaseStyle.Copy().Background(lipgloss.Color(GraidentColours[i])).PaddingLeft(2).Render
		} else {
			letterStyle = lettBaseStyle.Copy().Background(lipgloss.Color(GraidentColours[i])).Render
		}

		builder.WriteString(letterStyle(s))
		if i+1 == len(in) {
			builder.WriteString(letterStyle(" "))
		}
	}

	return "\n" + builder.String(), nil
}

func CopyDir(source string, destination string) error {
	// Create destination directory
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	// Get a list of files in the source directory
	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		sourcePath := filepath.Join(source, file.Name())
		destPath := filepath.Join(destination, file.Name())

		if file.IsDir() {
			// Recursively copy subdirectories
			err = CopyDir(sourcePath, destPath)
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err = copyFile(sourcePath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func CleanBuildDir(buildDir string) error {
	err := os.RemoveAll(buildDir)
	if err != nil {
		return err
	}
	os.Mkdir(buildDir, 0755)

	return nil
}

func ZipDirectory(source string, target string) (int64, error) {

	var zipFileSize int64

	zipFile, err := os.Create(target)
	if err != nil {
		return zipFileSize, err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	zipFileStat, err := zipFile.Stat()
	if err != nil {
		return zipFileSize, err
	}

	zipFileSize = zipFileStat.Size()

	return zipFileSize, nil
}

func Openbrowser(url string) error {
	var err error

	if RunningInDocker() {
		return nil
	}

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return err
	}

	return nil
}

func SectionMessage(message ...string) string {
	sectionOutputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderLeft(true).BorderBackground(lipgloss.Color(GraidentColours[0])).
		PaddingLeft(1).
		Render(strings.Join(message, " "))

	return sectionOutputStyle

}

func SectionErrorMessage(message ...string) string {
	sectionOutputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderLeft(true).BorderBackground(lipgloss.Color("#FF5658")).
		PaddingLeft(1).
		Render(Red(strings.Join(message, " ")))

	return sectionOutputStyle
}
