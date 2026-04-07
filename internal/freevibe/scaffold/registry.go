package scaffold

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed all:parenttmpl
var parentTmpl embed.FS

func loadRegistry(projectRoot string) (Registry, error) {
	path := filepath.Join(projectRoot, RegistryFileName)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Registry{Version: 1, Modules: []Module{}}, nil
		}
		return Registry{}, fmt.Errorf("read registry: %w", err)
	}
	var reg Registry
	if err := yaml.Unmarshal(b, &reg); err != nil {
		return Registry{}, fmt.Errorf("unmarshal registry: %w", err)
	}
	if reg.Version == 0 {
		reg.Version = 1
	}
	if reg.Modules == nil {
		reg.Modules = []Module{}
	}
	return reg, nil
}

func saveRegistry(projectRoot string, reg Registry) error {
	b, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}
	path := filepath.Join(projectRoot, RegistryFileName)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write registry: %w", err)
	}
	return nil
}

// copyParentTemplate copies all files from the embedded parenttmpl directory to
// projectRoot, skipping files that already exist. The __PROJECT_NAME__ token in
// file contents is replaced with projectName.
func copyParentTemplate(projectRoot, projectName string) error {
	return fs.WalkDir(parentTmpl, "parenttmpl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Compute destination path by stripping the "parenttmpl/" prefix.
		rel, err := filepath.Rel("parenttmpl", path)
		if err != nil {
			return err
		}
		dst := filepath.Join(projectRoot, rel)

		// Skip files that already exist.
		if _, err := os.Stat(dst); err == nil {
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
		}

		content, err := parentTmpl.ReadFile(path)
		if err != nil {
			return err
		}
		if projectName != "" {
			content = bytes.ReplaceAll(content, []byte("__PROJECT_NAME__"), []byte(projectName))
		}
		if err := os.WriteFile(dst, content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dst, err)
		}
		return nil
	})
}

func generateRootMakefile(projectRoot string) error {
	content := `.PHONY: modules status pull

modules:
	@cat freevibe.modules.yaml

status:
	@git submodule status

pull:
	@git submodule update --init --recursive
`
	path := filepath.Join(projectRoot, "Makefile")
	return os.WriteFile(path, []byte(content), 0o644)
}
