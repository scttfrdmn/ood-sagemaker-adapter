package cmd

import (
	"testing"

	smtypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

func TestAppReadyState(t *testing.T) {
	tests := []struct {
		status smtypes.AppStatus
		want   string
	}{
		{smtypes.AppStatusInService, "ready"},
		{smtypes.AppStatusPending, "pending"},
		{smtypes.AppStatusFailed, "failed"},
		{smtypes.AppStatusDeleting, "failed"},
		{smtypes.AppStatusDeleted, "failed"},
		{smtypes.AppStatus("SomeFutureStatus"), "pending"},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := appReadyState(tt.status); got != tt.want {
				t.Errorf("appReadyState(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}
