package scaffold

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	ModuleTypeBackend    = "backend"
	ModuleTypeFrontend   = "frontend"
	TemplateKratos       = "kratos"
	TemplateConsoleReact = "console-react"
	TemplatePythonTool   = "python-tool"

	RegistryFileName       = "freevibe.modules.yaml"
	DefaultTemplateRepoURL = "https://github.com/richer421/free-vibe.git"
)

var moduleNamePattern = regexp.MustCompile(`^[^\s/\\]+$`)

type TemplateSpec struct {
	Name              string
	Subdir            string
	DefaultModuleType string
}

var templateSpecs = map[string]TemplateSpec{
	TemplateKratos: {
		Name:              TemplateKratos,
		Subdir:            "templates/kratos",
		DefaultModuleType: ModuleTypeBackend,
	},
	TemplateConsoleReact: {
		Name:              TemplateConsoleReact,
		Subdir:            "templates/console-react",
		DefaultModuleType: ModuleTypeFrontend,
	},
	TemplatePythonTool: {
		Name:              TemplatePythonTool,
		Subdir:            "templates/python-tool",
		DefaultModuleType: TemplatePythonTool,
	},
}

type PromptFunc func(repoURL, defaultBranch string) (bool, error)

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
	ModuleName      string
	RepoURL         string
	Template        string
	TemplateRepoURL string
	Prompt          PromptFunc
}

type AddOptions struct {
	Name            string
	Type            string
	RepoURL         string
	ProjectName     string
	Template        string
	TemplateRepoURL string
	Prompt          PromptFunc
}

type renderData struct {
	ProjectName string
	ModuleName  string
}

// ValidModuleTypes returns the set of accepted module type strings, derived
// from the DefaultModuleType of every registered template.
func ValidModuleTypes() map[string]struct{} {
	types := make(map[string]struct{})
	for _, spec := range templateSpecs {
		types[spec.DefaultModuleType] = struct{}{}
	}
	return types
}

func ListTemplates() []string {
	names := make([]string, 0, len(templateSpecs))
	for name := range templateSpecs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func ResolveTemplate(name string) (TemplateSpec, error) {
	key := strings.TrimSpace(name)
	if key == "" {
		return TemplateSpec{}, fmt.Errorf("template is required")
	}
	spec, ok := templateSpecs[key]
	if !ok {
		return TemplateSpec{}, fmt.Errorf("unknown template: %s (run `freevibe template ls` to view available templates)", key)
	}
	return spec, nil
}

func ResolveTemplateSubdir(name string) (string, error) {
	spec, err := ResolveTemplate(name)
	if err != nil {
		return "", err
	}
	return spec.Subdir, nil
}
