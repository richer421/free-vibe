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
		Template:        TemplateKratos,
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

func TestInitProjectCreatesRootAgentsFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	remoteRepo := createBareRepo(t, filepath.Join(tmpDir, "ops-service.git"))
	setRemoteHead(t, remoteRepo, "main")
	projectRoot := filepath.Join(tmpDir, "ops-parent")

	err := InitProject(InitOptions{
		ProjectName:     "ops-parent",
		ProjectPath:     projectRoot,
		RepoURL:         remoteRepo,
		Template:        TemplateKratos,
		TemplateRepoURL: templateRepo,
	})
	if err != nil {
		t.Fatalf("InitProject error = %v", err)
	}

	agentsData, err := os.ReadFile(filepath.Join(projectRoot, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS: %v", err)
	}

	for _, snippet := range []string{
		"这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目",
		"这是父项目根仓库，不承载具体业务模块代码",
		"先查看 `freevibe.modules.yaml` 和 `.gitmodules`",
		"`knowledge/`",
		".codex/skills/codex-submodule-worktree-best-practices/SKILL.md",
		"进入对应子模块目录，并遵循该子模块自己的 `AGENTS.md` 与本地 `.codex` 约束",
	} {
		if !strings.Contains(string(agentsData), snippet) {
			t.Fatalf("AGENTS missing %q:\n%s", snippet, agentsData)
		}
	}
}

func TestInitProjectCreatesRootKnowledgeDirectory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	remoteRepo := createBareRepo(t, filepath.Join(tmpDir, "user-service.git"))
	setRemoteHead(t, remoteRepo, "main")
	projectRoot := filepath.Join(tmpDir, "user-parent")

	err := InitProject(InitOptions{
		ProjectName:     "user-parent",
		ProjectPath:     projectRoot,
		RepoURL:         remoteRepo,
		Template:        TemplateKratos,
		TemplateRepoURL: templateRepo,
	})
	if err != nil {
		t.Fatalf("InitProject error = %v", err)
	}

	knowledgeData, err := os.ReadFile(filepath.Join(projectRoot, "knowledge", "README.md"))
	if err != nil {
		t.Fatalf("read knowledge README: %v", err)
	}

	for _, snippet := range []string{
		"这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目",
		"根级 knowledge 用于沉淀父项目层面的业务背景、服务边界、跨模块协作约束和公共术语。",
	} {
		if !strings.Contains(string(knowledgeData), snippet) {
			t.Fatalf("knowledge README missing %q:\n%s", snippet, knowledgeData)
		}
	}
}

func TestInitProjectCreatesRootCodexSkill(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	templateRepo := createTemplateRepo(t, tmpDir)
	remoteRepo := createBareRepo(t, filepath.Join(tmpDir, "workflow-service.git"))
	setRemoteHead(t, remoteRepo, "main")
	projectRoot := filepath.Join(tmpDir, "workflow-parent")

	err := InitProject(InitOptions{
		ProjectName:     "workflow-parent",
		ProjectPath:     projectRoot,
		RepoURL:         remoteRepo,
		Template:        TemplateKratos,
		TemplateRepoURL: templateRepo,
	})
	if err != nil {
		t.Fatalf("InitProject error = %v", err)
	}

	skillData, err := os.ReadFile(filepath.Join(projectRoot, ".codex", "skills", "codex-submodule-worktree-best-practices", "SKILL.md"))
	if err != nil {
		t.Fatalf("read root codex skill: %v", err)
	}

	for _, snippet := range []string{
		"name: codex-submodule-worktree-best-practices",
		"项目关键词约定",
		"接收",
		"子模块指针更新",
	} {
		if !strings.Contains(string(skillData), snippet) {
			t.Fatalf("root codex skill missing %q:\n%s", snippet, skillData)
		}
	}
}

func TestRepositoryRootAgentsDescribesProject(t *testing.T) {
	t.Parallel()

	agentsPath := filepath.Join("..", "..", "..", "AGENTS.md")
	agentsData, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read repository root AGENTS: %v", err)
	}

	for _, snippet := range []string{
		"这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目",
		"目标不是堆功能，而是在满足业务目标的前提下，交付简洁、优雅、可扩展、可验证的结果。",
	} {
		if !strings.Contains(string(agentsData), snippet) {
			t.Fatalf("repository root AGENTS missing %q:\n%s", snippet, agentsData)
		}
	}
}

func TestRepositoryRootKnowledgeDirectoryExists(t *testing.T) {
	t.Parallel()

	knowledgePath := filepath.Join("..", "..", "..", "knowledge", "README.md")
	knowledgeData, err := os.ReadFile(knowledgePath)
	if err != nil {
		t.Fatalf("read repository root knowledge README: %v", err)
	}

	for _, snippet := range []string{
		"这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目",
		"根级 knowledge 用于沉淀父项目层面的业务背景、服务边界、跨模块协作约束和公共术语。",
	} {
		if !strings.Contains(string(knowledgeData), snippet) {
			t.Fatalf("repository root knowledge README missing %q:\n%s", snippet, knowledgeData)
		}
	}
}

func TestRepositoryRootCodexSkillExists(t *testing.T) {
	t.Parallel()

	skillPath := filepath.Join("..", "..", "..", ".codex", "skills", "codex-submodule-worktree-best-practices", "SKILL.md")
	skillData, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read repository root codex skill: %v", err)
	}

	for _, snippet := range []string{
		"name: codex-submodule-worktree-best-practices",
		"分支统一",
		"接收",
	} {
		if !strings.Contains(string(skillData), snippet) {
			t.Fatalf("repository root codex skill missing %q:\n%s", snippet, skillData)
		}
	}
}

func TestRepositoryRootAgentsReferencesCodexSubmoduleSkill(t *testing.T) {
	t.Parallel()

	agentsPath := filepath.Join("..", "..", "..", "AGENTS.md")
	agentsData, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read repository root AGENTS: %v", err)
	}

	if !strings.Contains(string(agentsData), ".codex/skills/codex-submodule-worktree-best-practices/SKILL.md") {
		t.Fatalf("repository root AGENTS missing codex submodule skill reference:\n%s", agentsData)
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
		Template:        TemplateKratos,
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
		Template:        TemplateKratos,
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

func TestKratosTemplateAgentsConnectsLocalSkills(t *testing.T) {
	t.Parallel()

	agentsPath := filepath.Join("..", "..", "..", "templates", "kratos", "AGENTS.md")
	agentsData, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read kratos AGENTS: %v", err)
	}

	for _, snippet := range []string{
		".codex/skills/free-vibe-coding-agent-core/SKILL.md",
		".codex/skills/business-domain-expert/SKILL.md",
		".codex/skills/kratos-layout-best-practices/SKILL.md",
	} {
		if !strings.Contains(string(agentsData), snippet) {
			t.Fatalf("kratos AGENTS missing %q:\n%s", snippet, agentsData)
		}
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
	if err := ensureProjectKnowledge(projectRoot); err != nil {
		t.Fatalf("ensure knowledge: %v", err)
	}
	if err := ensureProjectCodex(projectRoot); err != nil {
		t.Fatalf("ensure codex: %v", err)
	}
	if err := ensureProjectAgents(projectRoot); err != nil {
		t.Fatalf("ensure AGENTS: %v", err)
	}
	if err := generateRootMakefile(projectRoot); err != nil {
		t.Fatalf("generate Makefile: %v", err)
	}
	return projectRoot
}

func createTemplateRepo(t *testing.T, baseDir string) string {
	t.Helper()

	templateRoot := filepath.Join(baseDir, "free-vibe-coding-template")
	kratosRoot := filepath.Join(templateRoot, "templates", "kratos")
	if err := os.MkdirAll(templateRoot, 0o755); err != nil {
		t.Fatalf("mkdir template: %v", err)
	}
	runTestCmd(t, templateRoot, "git", "init", "-b", "main")
	runTestCmd(t, templateRoot, "git", "config", "user.email", "test@example.com")
	runTestCmd(t, templateRoot, "git", "config", "user.name", "test")
	if err := os.MkdirAll(filepath.Join(kratosRoot, "configs"), 0o755); err != nil {
		t.Fatalf("mkdir template subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(kratosRoot, "module.txt"), []byte("module=__MODULE_NAME__\nproject=__PROJECT_NAME__\n"), 0o644); err != nil {
		t.Fatalf("write template module.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(kratosRoot, "free-vibe-coding.txt"), []byte("free-vibe-coding\n"), 0o644); err != nil {
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
