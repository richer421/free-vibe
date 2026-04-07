package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

func renderModule(modulePath string, data renderData) error {
	replacements := map[string]string{
		"__PROJECT_NAME__": data.ProjectName,
		"__MODULE_NAME__":  data.ModuleName,
		"free-vibe-coding": data.ModuleName,
	}

	for _, token := range []string{"free-vibe-coding", "__MODULE_NAME__"} {
		if err := renamePathsWithToken(modulePath, token, data.ModuleName); err != nil {
			return err
		}
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

func copyTree(srcRoot, dstRoot string) error {
	return filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if rel == ".git" || strings.HasPrefix(rel, ".git"+string(filepath.Separator)) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dstPath := filepath.Join(dstRoot, rel)
		if d.IsDir() {
			return os.MkdirAll(dstPath, 0o755)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, b, 0o644)
	})
}

func clearWorktree(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(root, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}
