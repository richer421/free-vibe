package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const (
	ModuleTypeBackend  = "backend"
	ModuleTypeFrontend = "frontend"

	RegistryFileName       = "freevibe.modules.yaml"
	DefaultTemplateRepoURL = "https://github.com/richer421/free-vibe.git"
)

var moduleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,62}$`)

type Module struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Path    string `yaml:"path"`
	RepoURL string `yaml:"repo_url"`
}

type Registry struct {
	Version int      `yaml:"version"`
	Modules []Module `yaml:"modules"`
}

type InitOptions struct {
	ProjectName     string
	ProjectPath     string
	BackendName     string
	TemplateRepoURL string
	Render          bool
}

type AddOptions struct {
	Name        string
	Type        string
	RepoURL     string
	ProjectName string
	Render      bool
}

type renderData struct {
	ProjectName string
	ModuleName  string
}

func InitProject(opts InitOptions) error {
	if strings.TrimSpace(opts.ProjectName) == "" {
		return fmt.Errorf("project name is required")
	}
	if strings.TrimSpace(opts.ProjectPath) == "" {
		return fmt.Errorf("project path is required")
	}

	if err := os.MkdirAll(opts.ProjectPath, 0o755); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}
	if err := run(opts.ProjectPath, "git", "init"); err != nil {
		return fmt.Errorf("git init: %w", err)
	}

	registry := Registry{Version: 1, Modules: []Module{}}
	if err := saveRegistry(opts.ProjectPath, registry); err != nil {
		return err
	}
	if err := ensureProjectReadme(opts.ProjectPath, opts.ProjectName); err != nil {
		return err
	}

	backendName := strings.TrimSpace(opts.BackendName)
	if backendName != "" {
		repoURL := strings.TrimSpace(opts.TemplateRepoURL)
		if repoURL == "" {
			repoURL = DefaultTemplateRepoURL
		}
		if err := AddModule(opts.ProjectPath, AddOptions{
			Name:        backendName,
			Type:        ModuleTypeBackend,
			RepoURL:     repoURL,
			ProjectName: opts.ProjectName,
			Render:      opts.Render,
		}); err != nil {
			return fmt.Errorf("init backend module: %w", err)
		}
	}

	if err := generateRootMakefile(opts.ProjectPath); err != nil {
		return err
	}
	return nil
}

func AddModule(projectRoot string, opts AddOptions) error {
	moduleName := strings.TrimSpace(opts.Name)
	if err := validateModuleName(moduleName); err != nil {
		return err
	}
	moduleType := strings.TrimSpace(opts.Type)
	if moduleType != ModuleTypeBackend && moduleType != ModuleTypeFrontend {
		return fmt.Errorf("unsupported module type: %s", moduleType)
	}
	repoURL := strings.TrimSpace(opts.RepoURL)
	if repoURL == "" {
		return fmt.Errorf("repo URL is required")
	}

	registry, err := loadRegistry(projectRoot)
	if err != nil {
		return err
	}
	for _, m := range registry.Modules {
		if m.Name == moduleName || m.Path == moduleName {
			return fmt.Errorf("module already exists: %s", moduleName)
		}
	}

	if err := run(projectRoot, "git", "-c", "protocol.file.allow=always", "submodule", "add", repoURL, moduleName); err != nil {
		return fmt.Errorf("add submodule %s: %w", moduleName, err)
	}

	if opts.Render {
		rd := renderData{ProjectName: opts.ProjectName, ModuleName: moduleName}
		if err := renderModule(filepath.Join(projectRoot, moduleName), rd); err != nil {
			return fmt.Errorf("render module %s: %w", moduleName, err)
		}
	}

	registry.Modules = append(registry.Modules, Module{
		Name:    moduleName,
		Type:    moduleType,
		Path:    moduleName,
		RepoURL: repoURL,
	})
	if err := saveRegistry(projectRoot, registry); err != nil {
		return err
	}
	if err := generateRootMakefile(projectRoot); err != nil {
		return err
	}
	return nil
}

func RemoveModule(projectRoot, moduleName string) error {
	name := strings.TrimSpace(moduleName)
	if name == "" {
		return fmt.Errorf("module name is required")
	}

	registry, err := loadRegistry(projectRoot)
	if err != nil {
		return err
	}

	idx := -1
	modPath := ""
	for i := range registry.Modules {
		if registry.Modules[i].Name == name {
			idx = i
			modPath = registry.Modules[i].Path
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("module not found: %s", name)
	}

	if err := run(projectRoot, "git", "submodule", "deinit", "-f", "--", modPath); err != nil {
		return fmt.Errorf("deinit submodule %s: %w", modPath, err)
	}
	if err := run(projectRoot, "git", "rm", "-f", "--", modPath); err != nil {
		return fmt.Errorf("git rm submodule %s: %w", modPath, err)
	}

	_ = os.RemoveAll(filepath.Join(projectRoot, ".git", "modules", modPath))

	registry.Modules = append(registry.Modules[:idx], registry.Modules[idx+1:]...)
	if err := saveRegistry(projectRoot, registry); err != nil {
		return err
	}
	if err := generateRootMakefile(projectRoot); err != nil {
		return err
	}
	return nil
}

func validateModuleName(name string) error {
	if !moduleNamePattern.MatchString(name) {
		return fmt.Errorf("invalid module name: %s", name)
	}
	return nil
}

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
	content := fmt.Sprintf("# %s\n\nManaged by FreeVibe CLI.\n\n- Add module: `freevibe add --name <module>`\n- Remove module: `freevibe remove <module>`\n- Sync modules: `git submodule update --init --recursive`\n", projectName)
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

func renderModule(modulePath string, data renderData) error {
	replacements := map[string]string{
		"__PROJECT_NAME__": data.ProjectName,
		"__MODULE_NAME__":  data.ModuleName,
		"free-vibe-coding": data.ModuleName,
	}

	if err := renamePathsWithToken(modulePath, "free-vibe-coding", data.ModuleName); err != nil {
		return err
	}

	return filepath.WalkDir(modulePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !utf8.Valid(b) {
			return nil
		}

		updated := b
		for oldValue, newValue := range replacements {
			if oldValue == "" || oldValue == newValue {
				continue
			}
			updated = bytes.ReplaceAll(updated, []byte(oldValue), []byte(newValue))
		}
		if !bytes.Equal(b, updated) {
			if err := os.WriteFile(path, updated, 0o644); err != nil {
				return err
			}
		}
		return nil
	})
}

func renamePathsWithToken(root, token, replacement string) error {
	if token == "" || token == replacement {
		return nil
	}

	type renameItem struct {
		from string
		to   string
	}

	items := make([]renameItem, 0, 4)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, string(filepath.Separator)+".git") {
			return nil
		}
		base := filepath.Base(path)
		if !strings.Contains(base, token) {
			return nil
		}
		targetBase := strings.ReplaceAll(base, token, replacement)
		target := filepath.Join(filepath.Dir(path), targetBase)
		items = append(items, renameItem{from: path, to: target})
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(items, func(i, j int) bool {
		return len(items[i].from) > len(items[j].from)
	})

	for _, item := range items {
		if err := os.Rename(item.from, item.to); err != nil {
			return fmt.Errorf("rename %s -> %s: %w", item.from, item.to, err)
		}
	}
	return nil
}

func run(cwd, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return nil
}
