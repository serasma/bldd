package domain

import (
	"fmt"
	"os"
	"strings"

	"github.com/serrasma/bldd/internal/config"
)

type Reporter struct {
	cfg *config.Config
}

func NewReporter(cfg *config.Config) *Reporter {
	return &Reporter{
		cfg: cfg,
	}
}

func (r *Reporter) Render(librariesUsage []LibraryUsage) error {
	f, err := os.Create(r.cfg.ReportPath)
	if err != nil {
		return fmt.Errorf("reporter: create report file error: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("report on dynamic used libraries by ELF executables on %s\n\n",
		strings.Join(r.cfg.Directories, ",")))

	for _, libraryUsage := range librariesUsage {
		sb.WriteString(fmt.Sprintf("### %s\n", libraryUsage.Name))
		for _, file := range libraryUsage.Files {
			sb.WriteString(fmt.Sprintf("- %s\n", file))
		}
		sb.WriteString("\n")
	}

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("reporter: write report error: %w", err)
	}

	return nil
}
