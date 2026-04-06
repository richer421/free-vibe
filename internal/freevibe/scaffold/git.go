package scaffold

import (
	"fmt"
	"os/exec"
	"strings"
)

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

func runOutput(cwd, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return string(out), nil
}

func hasGitCommits(repoDir string) bool {
	return run(repoDir, "git", "rev-parse", "--verify", "HEAD") == nil
}

func hasGitChanges(repoDir string) (bool, error) {
	out, err := runOutput(repoDir, "git", "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}
