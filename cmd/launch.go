package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	smtypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/ood"
	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/sagemaker"
	"github.com/spf13/cobra"
)

var waitTimeout time.Duration

// appReadyState classifies a SageMaker AppStatus for the launch wait loop:
// "ready" once InService, "failed" on a terminal non-serviceable state, else
// "pending" (keep waiting). Pure function — unit-tested.
func appReadyState(status smtypes.AppStatus) string {
	switch status {
	case smtypes.AppStatusInService:
		return "ready"
	case smtypes.AppStatusFailed, smtypes.AppStatusDeleting, smtypes.AppStatusDeleted:
		return "failed"
	default: // Pending (and any unknown future value) — keep waiting
		return "pending"
	}
}

// waitForInService polls DescribeApp until the app is InService, errors on a
// terminal non-serviceable status, or times out.
func waitForInService(ctx context.Context, client *sagemaker.Client, did, profile, appName, aType string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		out, err := client.DescribeApp(ctx, did, profile, appName, aType)
		if err != nil {
			return err
		}
		switch appReadyState(out.Status) {
		case "ready":
			return nil
		case "failed":
			return fmt.Errorf("SageMaker app %q entered non-serviceable status %q", appName, out.Status)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s waiting for SageMaker app %q to reach InService (last status: %s)", timeout, appName, out.Status)
		}
		time.Sleep(5 * time.Second)
	}
}

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

		// Wait for the app to reach InService before generating the presigned URL.
		// CreatePresignedDomainUrl succeeds immediately, but a URL handed back before
		// the app is serviceable opens to a Studio spinner that times out (#2).
		if err := waitForInService(ctx, client, did, profile, appName, aType, waitTimeout); err != nil {
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
	launchCmd.Flags().DurationVar(&waitTimeout, "wait-timeout", 5*time.Minute,
		"how long to wait for the SageMaker app to reach InService before returning the presigned URL")
	rootCmd.AddCommand(launchCmd)
}
