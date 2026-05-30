package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	awstypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/ood"
	"github.com/scttfrdmn/ood-sagemaker-adapter/internal/sagemaker"
	"github.com/spf13/cobra"
)

// status <domain-id>/<user-profile>/<app-type>/<app-name>
var statusCmd = &cobra.Command{
	Use:   "status <app-id>",
	Short: "Get status of a SageMaker Studio app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "/", 4)
		if len(parts) != 4 {
			return fmt.Errorf("app-id must be <domain-id>/<user-profile>/<app-type>/<app-name>")
		}
		did, profile, aType, appName := parts[0], parts[1], parts[2], parts[3]

		ctx := context.Background()
		client, err := sagemaker.New(ctx, region)
		if err != nil {
			return err
		}

		out, err := client.DescribeApp(ctx, did, profile, appName, aType)
		if err != nil {
			return err
		}

		js := ood.AppStatus{
			ID:     args[0],
			Status: smStateToOod(out.Status),
		}
		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func smStateToOod(s awstypes.AppStatus) string {
	switch s {
	case awstypes.AppStatusPending:
		return ood.StatusQueued
	case awstypes.AppStatusInService:
		return ood.StatusRunning
	case awstypes.AppStatusDeleted:
		return ood.StatusCompleted
	case awstypes.AppStatusFailed:
		return ood.StatusFailed
	default:
		return ood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
