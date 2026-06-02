package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/sagemaker"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <app-id>",
	Short: "Delete a SageMaker Studio app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "/", 4)
		if len(parts) != 4 {
			return fmt.Errorf("app-id must be <domain-id>/<user-profile>/<app-type>/<app-name>")
		}
		did, profile, aType, appName := parts[0], parts[1], parts[2], parts[3]

		ctx := context.Background()
		client, err := sagemaker.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}
		if err := client.DeleteApp(ctx, did, profile, appName, aType); err != nil {
			return err
		}
		fmt.Printf("App %s deleted\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
