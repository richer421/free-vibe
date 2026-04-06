package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InitProject(opts InitOptions) error {
	if strings.TrimSpace(opts.ProjectName) == "" {
		return fmt.Errorf("project name is required")
	}
	if strings.TrimSpace(opts.ProjectPath) == "" {
		return fmt.Errorf("project path is required")
	}
	if strings.TrimSpace(opts.RepoURL) == "" {
		return fmt.Errorf("repo URL is required")
	}
	if strings.TrimSpace(opts.Template) == "" {
		return fmt.Errorf("template is required")
	}
	templateSpec, err := ResolveTemplate(opts.Template)
	if err != nil {
		return err
	}

	moduleName, err := resolveModuleName(opts.ModuleName, opts.RepoURL)
	if err != nil {
		return err
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
	if err := generateRootMakefile(opts.ProjectPath); err != nil {
		return err
	}

	return AddModule(opts.ProjectPath, AddOptions{
		Name:            moduleName,
		Type:            templateSpec.DefaultModuleType,
		RepoURL:         opts.RepoURL,
		ProjectName:     opts.ProjectName,
		Template:        opts.Template,
		TemplateRepoURL: opts.TemplateRepoURL,
		Prompt:          opts.Prompt,
	})
}

func AddModule(projectRoot string, opts AddOptions) error {
	moduleName, err := resolveModuleName(opts.Name, opts.RepoURL)
	if err != nil {
		return err
	}
	if err := validateModuleName(moduleName); err != nil {
		return err
	}

	moduleType := strings.TrimSpace(opts.Type)
	repoURL := strings.TrimSpace(opts.RepoURL)
	if repoURL == "" {
		return fmt.Errorf("repo URL is required")
	}
	if strings.TrimSpace(opts.Template) == "" {
		return fmt.Errorf("template is required")
	}
	templateSpec, err := ResolveTemplate(opts.Template)
	if err != nil {
		return err
	}
	if moduleType == "" {
		moduleType = templateSpec.DefaultModuleType
	}
	if moduleType != ModuleTypeBackend && moduleType != ModuleTypeFrontend {
		return fmt.Errorf("unsupported module type: %s", moduleType)
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

	if err := ensureModuleRepoReady(repoPreparation{
		RepoURL:         repoURL,
		Template:        opts.Template,
		TemplateRepoURL: opts.TemplateRepoURL,
		Data: renderData{
			ProjectName: opts.ProjectName,
			ModuleName:  moduleName,
		},
		Prompt: opts.Prompt,
	}); err != nil {
		return err
	}

	if err := run(projectRoot, "git", "-c", "protocol.file.allow=always", "submodule", "add", repoURL, moduleName); err != nil {
		return fmt.Errorf("add submodule %s: %w", moduleName, err)
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

func resolveModuleName(explicitName, repoURL string) (string, error) {
	name := strings.TrimSpace(explicitName)
	if name != "" {
		return name, nil
	}
	return ModuleNameFromRepoURL(repoURL)
}

func validateModuleName(name string) error {
	if !moduleNamePattern.MatchString(name) {
		return fmt.Errorf("invalid module name: %s", name)
	}
	return nil
}
