package cmd

import (
	"github.com/spf13/cobra"
)

var (
	region      string
	domainID    string
	userProfile string
	appType     string
)

var version = "dev" // overridden at release time via -ldflags -X .../cmd.version

var rootCmd = &cobra.Command{
	Version: version,
	Use:   "ood-sagemaker-adapter",
	Short: "OOD compute adapter for AWS SageMaker Studio",
	Long:  "Translates Open OnDemand interactive app requests to SageMaker Studio API calls.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&domainID, "domain-id", "", "SageMaker Domain ID")
	rootCmd.PersistentFlags().StringVar(&userProfile, "user-profile", "ood-default", "SageMaker user profile name")
	rootCmd.PersistentFlags().StringVar(&appType, "app-type", "JupyterServer", "SageMaker app type")
}
