package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

var (
	initBackendName string
	initTemplateURL string
	initNoRender    bool
	initForce       bool
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a submodule-based project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectPath, err := filepath.Abs(projectName)
		if err != nil {
			return err
		}

		if !initForce {
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
			BackendName:     initBackendName,
			TemplateRepoURL: initTemplateURL,
			Render:          !initNoRender,
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

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initBackendName, "backend-name", "free-vibe-backend", "Initial backend module name")
	initCmd.Flags().StringVar(&initTemplateURL, "template-repo-url", scaffold.DefaultTemplateRepoURL, "Remote template repository URL")
	initCmd.Flags().BoolVar(&initNoRender, "no-render", false, "Pull template only, skip render")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Continue even if target directory exists")
}
