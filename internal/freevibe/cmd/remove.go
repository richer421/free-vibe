package cmd

import (
	"fmt"
	"os"

	"free-vibe-coding/internal/freevibe/scaffold"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [module-name]",
	Short: "Remove a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, err := os.Getwd()
		if err != nil {
			return err
		}
		moduleName := args[0]
		fmt.Printf("Removing module: %s\n", moduleName)
		return scaffold.RemoveModule(projectRoot, moduleName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
