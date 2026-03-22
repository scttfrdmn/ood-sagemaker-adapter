package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/ood"
	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/sagemaker"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Create a SageMaker Studio app and return a presigned URL",
	Long:  "Reads a JSON app spec from stdin, creates the app, and prints the presigned Studio URL.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec ood.AppSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode app spec: %w", err)
		}

		did := domainID
		if did == "" {
			did = spec.DomainID
		}
		if did == "" {
			return fmt.Errorf("--domain-id is required")
		}

		ctx := context.Background()
		client, err := sagemaker.New(ctx, region)
		if err != nil {
			return err
		}

		appName := spec.AppName
		if appName == "" {
			appName = "ood-session"
		}
		profile := userProfile
		if spec.UserName != "" {
			profile = spec.UserName
		}
		aType := appType
		if spec.AppType != "" {
			aType = spec.AppType
		}

		if err := client.CreateApp(ctx, did, profile, appName, aType); err != nil {
			return err
		}

		url, err := client.CreatePresignedURL(ctx, did, profile)
		if err != nil {
			return err
		}

		result := ood.AppStatus{
			ID:         fmt.Sprintf("%s/%s/%s/%s", did, profile, aType, appName),
			Status:     ood.StatusRunning,
			PresignURL: url,
		}
		return json.NewEncoder(os.Stdout).Encode(result)
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
