package cmd

import (
	"fmt"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Template commands",
	}
	cmd.AddCommand(newTemplateLsCmd())
	return cmd
}

func newTemplateLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List available templates",
		Run: func(cmd *cobra.Command, args []string) {
			for _, name := range scaffold.ListTemplates() {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
		},
	}
}
