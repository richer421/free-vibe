package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

var (
	addName     string
	addType     string
	addRepoURL  string
	addNoRender bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a module",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, err := os.Getwd()
		if err != nil {
			return err
		}
		projectName := filepath.Base(projectRoot)

		if addName == "" {
			return fmt.Errorf("--name is required")
		}

		repoURL := addRepoURL
		if repoURL == "" {
			repoURL = scaffold.DefaultTemplateRepoURL
		}

		fmt.Printf("Adding module: %s (%s)\n", addName, addType)
		return scaffold.AddModule(projectRoot, scaffold.AddOptions{
			Name:        addName,
			Type:        addType,
			RepoURL:     repoURL,
			ProjectName: projectName,
			Render:      !addNoRender,
		})
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addName, "name", "", "Module name, for example: order-service")
	addCmd.Flags().StringVar(&addType, "type", scaffold.ModuleTypeBackend, "Module type: backend/frontend")
	addCmd.Flags().StringVar(&addRepoURL, "repo-url", "", "Module repository URL (default: FreeVibe template repo)")
	addCmd.Flags().BoolVar(&addNoRender, "no-render", false, "Pull module only, skip render")
	_ = addCmd.MarkFlagRequired("name")
}
