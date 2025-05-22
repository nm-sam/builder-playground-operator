package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	controller "github.com/flashbots/builder-playground-operator/internal/controller"
)

var (
	cliMode          bool
	manifestPath     string
	outputDir        string
	builderConfigDir string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "builder-playground-operator",
		Short: "Builder Playground Operator CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cliMode {
				if manifestPath == "" || outputDir == "" {
					return fmt.Errorf("both --manifest and --k8s-manifests-dir are required when using --cli")
				}
				if builderConfigDir == "" {
					builderConfigDir = "$HOME/.playground/devnet"
				}
				return controller.GenerateCRAndStatefulSet(manifestPath, outputDir, builderConfigDir)
			}
			fmt.Println("Please specify a mode. Use --cli for CLI mode.")
			return nil
		},
	}

	// CLI flags
	rootCmd.Flags().BoolVar(&cliMode, "cli", false, "Run in CLI mode")
	rootCmd.Flags().StringVar(&manifestPath, "manifest", "", "Path to the manifest file (e.g., ./manifest.json)")
	rootCmd.Flags().StringVar(&outputDir, "k8s-manifests-dir", "", "Directory to write Kubernetes manifests to")
	rootCmd.Flags().StringVar(&builderConfigDir, "builder-config-dir", "", "Path to the builder playground configuration files like secrets")

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
