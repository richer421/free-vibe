package scaffold

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestModuleNameFromRepoURL(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"https://github.com/richer421/free-vibe.git": "free-vibe",
		"git@github.com:richer421/free-vibe.git":     "free-vibe",
		"/tmp/payment-service.git":                   "payment-service",
	}

	for repoURL, want := range cases {
		repoURL := repoURL
		want := want
		t.Run(repoURL, func(t *testing.T) {
			t.Parallel()

			got, err := ModuleNameFromRepoURL(repoURL)
			if err != nil {
				t.Fatalf("ModuleNameFromRepoURL(%q) error = %v", repoURL, err)
			}
			if got != want {
				t.Fatalf("ModuleNameFromRepoURL(%q) = %q, want %q", repoURL, got, want)
			}
		})
	}
}

func TestPromptBootstrapConfirmation(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	ok, err := PromptBootstrapConfirmation(
		strings.NewReader("yes\n"),
		&out,
		"https://github.com/richer421/existing-service.git",
		"main",
	)
	if err != nil {
		t.Fatalf("PromptBootstrapConfirmation error = %v", err)
	}
	if !ok {
		t.Fatal("PromptBootstrapConfirmation returned false, want true")
	}

	prompt := out.String()
	for _, snippet := range []string{
		"仓库：https://github.com/richer421/existing-service.git",
		"默认分支：main",
		"是否要用当前脚手架内容初始化这个仓库？",
		"输入 yes：",
		"输入 no：",
	} {
		if !strings.Contains(prompt, snippet) {
			t.Fatalf("prompt missing %q:\n%s", snippet, prompt)
		}
	}
}

func TestInitProjectBootstrapsEmptyRepoUsingDerivedModuleName(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	remoteRepo := createBareRepo(t, filepath.Join(tmpDir, "empty-service.git"))
	setRemoteHead(t, remoteRepo, "main")
	projectRoot := filepath.Join(tmpDir, "demo-monorepo")

	err := InitProject(InitOptions{
		ProjectName:     "demo-monorepo",
		ProjectPath:     projectRoot,
		RepoURL:         remoteRepo,
		TemplateRepoURL: templateRepo,
	})
	if err != nil {
		t.Fatalf("InitProject error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectRoot, "empty-service")); err != nil {
		t.Fatalf("expected derived module directory: %v", err)
	}

	registryData, err := os.ReadFile(filepath.Join(projectRoot, RegistryFileName))
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	if !strings.Contains(string(registryData), "name: empty-service") {
		t.Fatalf("registry missing derived module name:\n%s", registryData)
	}
	if !strings.Contains(string(registryData), "repo_url: "+remoteRepo) {
		t.Fatalf("registry missing target repo URL:\n%s", registryData)
	}

	cloneDir := filepath.Join(tmpDir, "verify-empty-service")
	runTestCmd(t, "", "git", "clone", remoteRepo, cloneDir)
	content, err := os.ReadFile(filepath.Join(cloneDir, "module.txt"))
	if err != nil {
		t.Fatalf("read bootstrapped file: %v", err)
	}
	if string(content) != "module=empty-service\nproject=demo-monorepo\n" {
		t.Fatalf("unexpected bootstrapped content: %q", string(content))
	}
}

func TestAddModuleAttachesExistingRepoWhenUserAnswersNo(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	projectRoot := createParentProject(t, tmpDir, "attach-parent")
	remoteRepo := createNonEmptyRemoteRepo(t, tmpDir, "inventory-service.git", map[string]string{
		"README.md": "existing content\n",
	})

	err := AddModule(projectRoot, AddOptions{
		ProjectName:     "attach-parent",
		Type:            ModuleTypeBackend,
		RepoURL:         remoteRepo,
		TemplateRepoURL: templateRepo,
		Prompt: func(repoURL, branch string) (bool, error) {
			if repoURL != remoteRepo {
				t.Fatalf("Prompt repoURL = %q, want %q", repoURL, remoteRepo)
			}
			if branch != "main" {
				t.Fatalf("Prompt branch = %q, want main", branch)
			}
			return false, nil
		},
	})
	if err != nil {
		t.Fatalf("AddModule error = %v", err)
	}

	registryData, err := os.ReadFile(filepath.Join(projectRoot, RegistryFileName))
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	if !strings.Contains(string(registryData), "name: inventory-service") {
		t.Fatalf("registry missing attached module:\n%s", registryData)
	}

	cloneDir := filepath.Join(tmpDir, "verify-attach")
	runTestCmd(t, "", "git", "clone", remoteRepo, cloneDir)
	content, err := os.ReadFile(filepath.Join(cloneDir, "README.md"))
	if err != nil {
		t.Fatalf("read attached remote content: %v", err)
	}
	if string(content) != "existing content\n" {
		t.Fatalf("existing remote content changed unexpectedly: %q", string(content))
	}
}

func TestAddModuleBootstrapsExistingRepoWhenUserAnswersYes(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	projectRoot := createParentProject(t, tmpDir, "bootstrap-parent")
	remoteRepo := createNonEmptyRemoteRepo(t, tmpDir, "payment-service.git", map[string]string{
		"README.md": "legacy content\n",
	})

	err := AddModule(projectRoot, AddOptions{
		ProjectName:     "bootstrap-parent",
		Type:            ModuleTypeBackend,
		RepoURL:         remoteRepo,
		TemplateRepoURL: templateRepo,
		Prompt: func(repoURL, branch string) (bool, error) {
			return true, nil
		},
	})
	if err != nil {
		t.Fatalf("AddModule error = %v", err)
	}

	cloneDir := filepath.Join(tmpDir, "verify-bootstrap")
	runTestCmd(t, "", "git", "clone", remoteRepo, cloneDir)
	content, err := os.ReadFile(filepath.Join(cloneDir, "module.txt"))
	if err != nil {
		t.Fatalf("read bootstrapped file: %v", err)
	}
	if string(content) != "module=payment-service\nproject=bootstrap-parent\n" {
		t.Fatalf("unexpected bootstrapped content: %q", string(content))
	}

	count := strings.TrimSpace(runTestCmd(t, cloneDir, "git", "rev-list", "--count", "HEAD"))
	if count != "2" {
		t.Fatalf("commit count = %s, want 2", count)
	}
}

func createParentProject(t *testing.T, baseDir, name string) string {
	t.Helper()

	projectRoot := filepath.Join(baseDir, name)
	if err := os.MkdirAll(projectRoot, 0o755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}
	runTestCmd(t, projectRoot, "git", "init", "-b", "main")
	if err := saveRegistry(projectRoot, Registry{Version: 1, Modules: []Module{}}); err != nil {
		t.Fatalf("save registry: %v", err)
	}
	if err := ensureProjectReadme(projectRoot, name); err != nil {
		t.Fatalf("ensure README: %v", err)
	}
	if err := generateRootMakefile(projectRoot); err != nil {
		t.Fatalf("generate Makefile: %v", err)
	}
	return projectRoot
}

func createTemplateRepo(t *testing.T, baseDir string) string {
	t.Helper()

	templateRoot := filepath.Join(baseDir, "free-vibe-coding-template")
	if err := os.MkdirAll(templateRoot, 0o755); err != nil {
		t.Fatalf("mkdir template: %v", err)
	}
	runTestCmd(t, templateRoot, "git", "init", "-b", "main")
	runTestCmd(t, templateRoot, "git", "config", "user.email", "test@example.com")
	runTestCmd(t, templateRoot, "git", "config", "user.name", "test")
	if err := os.MkdirAll(filepath.Join(templateRoot, "configs"), 0o755); err != nil {
		t.Fatalf("mkdir configs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateRoot, "module.txt"), []byte("module=__MODULE_NAME__\nproject=__PROJECT_NAME__\n"), 0o644); err != nil {
		t.Fatalf("write template module.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateRoot, "free-vibe-coding.txt"), []byte("free-vibe-coding\n"), 0o644); err != nil {
		t.Fatalf("write template renamed file: %v", err)
	}
	runTestCmd(t, templateRoot, "git", "add", ".")
	runTestCmd(t, templateRoot, "git", "commit", "-m", "template init")
	return templateRoot
}

func createBareRepo(t *testing.T, path string) string {
	t.Helper()
	runTestCmd(t, "", "git", "init", "--bare", path)
	return path
}

func setRemoteHead(t *testing.T, repoPath, branch string) {
	t.Helper()
	runTestCmd(t, "", "git", "--git-dir", repoPath, "symbolic-ref", "HEAD", "refs/heads/"+branch)
}

func createNonEmptyRemoteRepo(t *testing.T, baseDir, name string, files map[string]string) string {
	t.Helper()

	remoteRepo := createBareRepo(t, filepath.Join(baseDir, name))
	setRemoteHead(t, remoteRepo, "main")

	worktree := filepath.Join(baseDir, strings.TrimSuffix(name, ".git")+"-src")
	runTestCmd(t, "", "git", "clone", remoteRepo, worktree)
	runTestCmd(t, worktree, "git", "config", "user.email", "test@example.com")
	runTestCmd(t, worktree, "git", "config", "user.name", "test")

	for path, content := range files {
		fullPath := filepath.Join(worktree, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("mkdir file parent: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write file %s: %v", path, err)
		}
	}

	runTestCmd(t, worktree, "git", "add", ".")
	runTestCmd(t, worktree, "git", "commit", "-m", "initial content")
	runTestCmd(t, worktree, "git", "push", "origin", "HEAD:main")
	return remoteRepo
}

func runTestCmd(t *testing.T, cwd string, name string, args ...string) string {
	t.Helper()

	cmd := exec.Command(name, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return string(out)
}
