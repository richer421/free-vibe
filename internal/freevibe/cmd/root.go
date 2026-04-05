package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildTime = ""
)

var rootCmd = &cobra.Command{
	Use:   "freevibe",
	Short: "FreeVibe mono-project scaffold CLI",
	Long: `freevibe creates a parent repository with git submodules.

Core flow:
  1) freevibe init   create parent repo and bootstrap first module
  2) freevibe add    add one module as submodule
  3) freevibe remove remove one module`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("freevibe version %s", Version)
		if BuildTime != "" {
			fmt.Printf(" (built at %s)", BuildTime)
		}
		fmt.Println()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
