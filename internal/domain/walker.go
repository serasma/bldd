package domain

import (
	"context"
	"debug/elf"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/facebookgo/symwalk"
	"github.com/serrasma/bldd/internal/config"
	"golang.org/x/sync/errgroup"
)

type Walker struct {
	librariesUsage map[string][]string
	librariesList  map[string]struct{}
	directories    []string
	cfg            *config.Config
}

func NewWalker(cfg *config.Config) *Walker {
	librariesList := make(map[string]struct{})
	for _, library := range cfg.Libraries {
		librariesList[library] = struct{}{}
	}

	return &Walker{
		librariesList: librariesList,
		directories:   resolveDirectories(cfg.Directories),
		cfg:           cfg,
	}
}

func (w *Walker) Walk(ctx context.Context) ([]LibraryUsage, error) {
	eg, groupCtx := errgroup.WithContext(ctx)
	eg.SetLimit(w.cfg.Workers)

	for _, directory := range w.directories {
		eg.Go(func() error {
			if err := symwalk.Walk(directory, func(path string, info fs.FileInfo, err error) error {
				select {
				case <-groupCtx.Done():
					return context.Cause(groupCtx)
				default:
				}

				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				if err = w.listDependencies(path); err != nil {
					return err
				}

				return nil
			}); err != nil {
				return fmt.Errorf("walkdir error: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("walker: errgroup wait error: %w", err)
	}

	return w.usageSort(), nil
}

func (w *Walker) listDependencies(path string) error {
	f, err := elf.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	needs, err := f.DynamicVersionNeeds()
	if err != nil {
		return err
	}

	for _, need := range needs {
		if _, ok := w.librariesList[need.Name]; !ok {
			continue
		}

		libraryName := fmt.Sprintf("%s [%s]", need.Name, strings.TrimPrefix(f.Machine.String(), ELFMachine))
		w.librariesUsage[libraryName] = append(w.librariesUsage[libraryName], path)
	}

	return nil
}

func (w *Walker) usageSort() []LibraryUsage {
	librariesUsage := make([]LibraryUsage, 0, len(w.librariesUsage))
	for name, files := range w.librariesUsage {
		librariesUsage = append(librariesUsage, LibraryUsage{
			Name:  name,
			Files: files,
		})
	}

	slices.SortFunc(librariesUsage, func(lhs LibraryUsage, rhs LibraryUsage) int {
		if len(lhs.Files) < len(rhs.Files) {
			return -1
		}

		if len(lhs.Files) > len(rhs.Files) {
			return 1
		}

		if lhs.Name < rhs.Name {
			return -1
		}

		if lhs.Name > rhs.Name {
			return 1
		}

		return 0
	})

	return librariesUsage
}

func resolveDirectories(directories []string) []string {
	resolved := make([]string, 0, len(directories))
	for _, directory := range directories {
		target, err := filepath.EvalSymlinks(directory)
		if err == nil {
			resolved = append(resolved, target)
		} else {
			resolved = append(resolved, directory)
		}
	}

	return resolved
}

const (
	ELFMachine = "EM_"
)
