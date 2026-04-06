package scaffold

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type repoPreparation struct {
	RepoURL         string
	TemplateRepoURL string
	Data            renderData
	Prompt          PromptFunc
}

func ModuleNameFromRepoURL(repoURL string) (string, error) {
	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return "", fmt.Errorf("repo URL is required")
	}

	repoURL = strings.TrimSuffix(repoURL, "/")
	base := repoURL
	if strings.Contains(repoURL, "://") {
		base = repoURL[strings.LastIndex(repoURL, "/")+1:]
	} else if idx := strings.LastIndexAny(repoURL, "/:"); idx >= 0 {
		base = repoURL[idx+1:]
	}

	base = strings.TrimSuffix(base, ".git")
	base = filepath.Base(base)
	if strings.TrimSpace(base) == "" || base == "." {
		return "", fmt.Errorf("invalid repo URL: %s", repoURL)
	}
	return base, nil
}

func NewConsolePrompt(in io.Reader, out io.Writer) PromptFunc {
	return func(repoURL, defaultBranch string) (bool, error) {
		return PromptBootstrapConfirmation(in, out, repoURL, defaultBranch)
	}
}

func PromptBootstrapConfirmation(in io.Reader, out io.Writer, repoURL, defaultBranch string) (bool, error) {
	if in == nil {
		in = strings.NewReader("")
	}
	if out == nil {
		out = io.Discard
	}

	_, _ = fmt.Fprintf(out, "仓库：%s\n默认分支：%s\n\n是否要用当前脚手架内容初始化这个仓库？\n\n输入 yes：\n- 保留现有 git 历史\n- 用脚手架内容覆盖当前文件内容\n- 在默认分支新增一次提交\n\n输入 no：\n- 不写入脚手架\n- 直接把现有仓库作为子模块接入父项目\n", repoURL, defaultBranch)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch answer {
		case "yes":
			return true, nil
		case "no":
			return false, nil
		case "":
			continue
		default:
			_, _ = fmt.Fprintln(out, "请输入 yes 或 no")
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, fmt.Errorf("expected yes or no")
}

func ensureModuleRepoReady(opts repoPreparation) error {
	templateRepoURL := strings.TrimSpace(opts.TemplateRepoURL)
	if templateRepoURL == "" {
		templateRepoURL = DefaultTemplateRepoURL
	}

	repoDir, err := os.MkdirTemp("", "freevibe-module-repo-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(repoDir)

	if err := run("", "git", "clone", opts.RepoURL, repoDir); err != nil {
		return fmt.Errorf("clone repo %s: %w", opts.RepoURL, err)
	}

	defaultBranch, err := detectDefaultBranch(repoDir)
	if err != nil {
		return err
	}
	hasCommits := hasGitCommits(repoDir)

	if hasCommits {
		if opts.Prompt == nil {
			return fmt.Errorf("prompt is required for non-empty repository: %s", opts.RepoURL)
		}
		useScaffold, err := opts.Prompt(opts.RepoURL, defaultBranch)
		if err != nil {
			return err
		}
		if !useScaffold {
			return nil
		}
	}

	return bootstrapRepo(repoDir, defaultBranch, hasCommits, templateRepoURL, opts.Data)
}

func detectDefaultBranch(repoDir string) (string, error) {
	out, err := runOutput(repoDir, "git", "remote", "show", "origin")
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "HEAD branch:") {
			continue
		}
		branch := strings.TrimSpace(strings.TrimPrefix(line, "HEAD branch:"))
		if branch != "" && branch != "(unknown)" {
			return branch, nil
		}
	}

	currentBranch, err := runOutput(repoDir, "git", "branch", "--show-current")
	if err == nil && strings.TrimSpace(currentBranch) != "" {
		return strings.TrimSpace(currentBranch), nil
	}
	return "main", nil
}

func bootstrapRepo(repoDir, defaultBranch string, hasCommits bool, templateRepoURL string, data renderData) error {
	if err := checkoutBranch(repoDir, defaultBranch, hasCommits); err != nil {
		return err
	}
	if err := clearWorktree(repoDir); err != nil {
		return err
	}

	templateDir, err := os.MkdirTemp("", "freevibe-template-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(templateDir)

	if err := run("", "git", "clone", templateRepoURL, templateDir); err != nil {
		return fmt.Errorf("clone template repo %s: %w", templateRepoURL, err)
	}
	if err := renderModule(templateDir, data); err != nil {
		return err
	}
	if err := copyTree(templateDir, repoDir); err != nil {
		return err
	}

	if err := run(repoDir, "git", "add", "-A"); err != nil {
		return err
	}

	changed, err := hasGitChanges(repoDir)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	if err := run(repoDir, "git", "-c", "user.name=FreeVibe", "-c", "user.email=freevibe@local", "commit", "-m", fmt.Sprintf("chore: bootstrap %s with FreeVibe", data.ModuleName)); err != nil {
		return err
	}
	if err := run(repoDir, "git", "push", "origin", "HEAD:"+defaultBranch); err != nil {
		return err
	}
	return nil
}

func checkoutBranch(repoDir, defaultBranch string, hasCommits bool) error {
	if hasCommits {
		currentBranch, err := runOutput(repoDir, "git", "branch", "--show-current")
		if err == nil && strings.TrimSpace(currentBranch) == defaultBranch {
			return nil
		}
		if err := run(repoDir, "git", "checkout", defaultBranch); err == nil {
			return nil
		}
		return run(repoDir, "git", "checkout", "-B", defaultBranch, "origin/"+defaultBranch)
	}

	currentBranch, err := runOutput(repoDir, "git", "branch", "--show-current")
	if err == nil && strings.TrimSpace(currentBranch) == defaultBranch {
		return nil
	}
	if err := run(repoDir, "git", "checkout", "--orphan", defaultBranch); err == nil {
		return nil
	}
	return run(repoDir, "git", "symbolic-ref", "HEAD", "refs/heads/"+defaultBranch)
}
