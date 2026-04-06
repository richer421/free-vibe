package scaffold

import "regexp"

const (
	ModuleTypeBackend  = "backend"
	ModuleTypeFrontend = "frontend"

	RegistryFileName       = "freevibe.modules.yaml"
	DefaultTemplateRepoURL = "https://github.com/richer421/free-vibe.git"
	DefaultTemplateSubdir  = "templates/kratos"
)

var moduleNamePattern = regexp.MustCompile(`^[^\s/\\]+$`)

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
	TemplateRepoURL string
	Prompt          PromptFunc
}

type AddOptions struct {
	Name            string
	Type            string
	RepoURL         string
	ProjectName     string
	TemplateRepoURL string
	Prompt          PromptFunc
}

type renderData struct {
	ProjectName string
	ModuleName  string
}
