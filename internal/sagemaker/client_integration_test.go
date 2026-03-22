//go:build integration

package sagemaker_test

import (
	"context"
	"strings"
	"testing"

	substrate "github.com/scttfrdmn/substrate"

	. "github.com/scttfrdmn/ood-sagemaker-adapter/internal/sagemaker"
)

const (
	testDomainID    = "d-testdomain0001"
	testUserProfile = "test-user"
)

// TestCreateDescribeDeleteApp_Substrate exercises the full SageMaker Studio app
// lifecycle (CreateApp → DescribeApp → DeleteApp) against the substrate emulator.
func TestCreateDescribeDeleteApp_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	appName := "ood-test-app"
	appType := "JupyterServer"

	// CreateApp
	err = client.CreateApp(ctx, testDomainID, testUserProfile, appName, appType)
	if err != nil {
		t.Fatalf("CreateApp: %v", err)
	}
	t.Log("app created successfully")

	// DescribeApp
	detail, err := client.DescribeApp(ctx, testDomainID, testUserProfile, appName, appType)
	if err != nil {
		t.Fatalf("DescribeApp: %v", err)
	}
	if detail == nil {
		t.Fatal("DescribeApp: got nil output")
	}
	t.Logf("app status: %s", detail.Status)

	// DeleteApp
	err = client.DeleteApp(ctx, testDomainID, testUserProfile, appName, appType)
	if err != nil {
		t.Fatalf("DeleteApp: %v", err)
	}
	t.Log("app deleted successfully")
}

// TestDescribeApp_NotFound_Substrate verifies that DescribeApp returns an error
// for an app that was never created.
func TestDescribeApp_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.DescribeApp(ctx, "d-doesnotexist", "nobody", "no-app", "JupyterServer")
	if err == nil {
		t.Fatal("expected error for non-existent app, got nil")
	}
	if !strings.Contains(err.Error(), "sagemaker") {
		t.Logf("error (acceptable): %v", err)
	}
}
