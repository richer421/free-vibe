package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var name string
	var moduleType string
	var repoURL string
	var template string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a module",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, err := os.Getwd()
			if err != nil {
				return err
			}
			projectName := filepath.Base(projectRoot)

			fmt.Printf("Adding module from repo: %s\n", repoURL)
			return scaffold.AddModule(projectRoot, scaffold.AddOptions{
				Name:        name,
				Type:        moduleType,
				RepoURL:     repoURL,
				ProjectName: projectName,
				Template:    template,
				Prompt:      scaffold.NewConsolePrompt(cmd.InOrStdin(), cmd.OutOrStdout()),
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Module name (defaults to repo name)")
	cmd.Flags().StringVar(&moduleType, "type", "", "Module type (optional: backend/frontend, defaults to template type)")
	cmd.Flags().StringVar(&repoURL, "repo", "", "Target module repository URL")
	cmd.Flags().StringVar(&template, "template", "", "Template name")
	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("template")
	return cmd
}
