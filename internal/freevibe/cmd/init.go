package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var moduleName string
	var repoURL string
	var template string
	var templateRepo string
	var force bool

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a submodule-based project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			projectPath, err := filepath.Abs(projectName)
			if err != nil {
				return err
			}

			if !force {
				if st, err := os.Stat(projectPath); err == nil && st.IsDir() {
					entries, readErr := os.ReadDir(projectPath)
					if readErr != nil {
						return readErr
					}
					if len(entries) > 0 {
						return fmt.Errorf("target dir is not empty: %s (use --force to continue)", projectPath)
					}
				}
			}

			fmt.Printf("Initializing FreeVibe project: %s\n", projectName)
			err = scaffold.InitProject(scaffold.InitOptions{
				ProjectName:     projectName,
				ProjectPath:     projectPath,
				ModuleName:      moduleName,
				RepoURL:         repoURL,
				Template:        template,
				TemplateRepoURL: templateRepo,
				Prompt:          scaffold.NewConsolePrompt(cmd.InOrStdin(), cmd.OutOrStdout()),
			})
			if err != nil {
				return err
			}

			fmt.Println("Project initialized")
			fmt.Printf("Path: %s\n", projectPath)
			fmt.Println("Next: cd <project> && git submodule status")
			return nil
		},
	}

	cmd.Flags().StringVar(&moduleName, "name", "", "Initial module name (defaults to repo name)")
	cmd.Flags().StringVar(&repoURL, "repo", "", "Target module repository URL")
	cmd.Flags().StringVar(&template, "template", "", "Template name")
	cmd.Flags().StringVar(&templateRepo, "template-repo", "", "Template repository URL (defaults to official repo)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Continue even if target directory exists")
	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("template")
	return cmd
}
