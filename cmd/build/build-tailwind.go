package build

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/damongolding/eventsforce/internal/utils"
)

type TailwindError struct {
	Reason string
	File   string
	Line   string
}

func (t *TailwindError) parseErrors(buff *bytes.Buffer) {
	scanner := bufio.NewScanner(buff)
	for scanner.Scan() {
		switch {
		case strings.Contains(scanner.Text(), "reason: '"):
			theReason := strings.Split(scanner.Text(), "'")
			t.Reason = theReason[1]
		case strings.Contains(scanner.Text(), "file: '"):
			theFile := strings.Split(scanner.Text(), "'")
			t.File = theFile[1]
		case strings.Contains(scanner.Text(), "line: "):
			theLine := strings.Split(scanner.Text(), ":")
			t.File = strings.TrimSpace(theLine[1])
		}
	}
}

func (t *TailwindError) getError() error {

	heading := utils.Red(lipgloss.NewStyle().Bold(true).Render("Tailwind error"))
	reason := utils.Red(lipgloss.NewStyle().Render("   Reason :", t.Reason))
	file := utils.Red(lipgloss.NewStyle().Render("   File   :", utils.RemoveDockerPathPrefix(t.File)))
	line := utils.Red(lipgloss.NewStyle().Render("   Line   :", t.Line))

	out := lipgloss.JoinVertical(lipgloss.Left, heading, reason, file, line)

	return fmt.Errorf("%s", out)
}

func Tailwind(ctx context.Context, path string, args ...string) error {

	htmlPath := filepath.Dir(path)
	htmlPath = fmt.Sprintf("%s/index.html", htmlPath)
	configPath := "./tailwind.config.js"
	if config.InContainer {
		configPath = "/templates/tailwind.config.js"
	}
	args = slices.Concat([]string{"--content", htmlPath, "--config", configPath, "-i", path, "-o", path}, args)

	cmd := exec.CommandContext(ctx, "./tailwindcss", args...)

	s := ""
	buff := bytes.NewBufferString(s)
	cmd.Stderr = buff

	if err := cmd.Run(); err != nil {
		var tailwindErr TailwindError
		tailwindErr.parseErrors(buff)
		return tailwindErr.getError()
	}

	return nil
}
