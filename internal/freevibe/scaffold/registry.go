package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

func ensureProjectReadme(projectRoot, projectName string) error {
	path := filepath.Join(projectRoot, "README.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	content := fmt.Sprintf("# %s\n\nManaged by FreeVibe CLI.\n\n- Add module: `freevibe add --repo <url>`\n- Remove module: `freevibe remove <module>`\n- Sync modules: `git submodule update --init --recursive`\n", projectName)
	return os.WriteFile(path, []byte(content), 0o644)
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
