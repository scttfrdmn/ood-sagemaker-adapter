// Package sagemaker wraps the AWS SageMaker API for the OOD adapter.
package sagemaker

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

// Client wraps the AWS SageMaker client.
type Client struct {
	svc    *sagemaker.Client
	region string
}

// New creates a SageMaker client using the default AWS credential chain.
func New(ctx context.Context, region string, optFns ...func(*config.LoadOptions) error) (*Client, error) {
	if len(optFns) == 0 {
		optFns = []func(*config.LoadOptions) error{config.WithRegion(region)}
	}
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: sagemaker.NewFromConfig(cfg), region: region}, nil
}

// CreateApp creates a SageMaker app for the given user profile.
func (c *Client) CreateApp(ctx context.Context, domainID, userProfile, appName, appType string) error {
	_, err := c.svc.CreateApp(ctx, &sagemaker.CreateAppInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfile),
		AppName:         aws.String(appName),
		AppType:         types.AppType(appType),
	})
	if err != nil {
		return fmt.Errorf("sagemaker CreateApp: %w", err)
	}
	return nil
}

// DescribeApp returns the current status of a SageMaker app.
func (c *Client) DescribeApp(ctx context.Context, domainID, userProfile, appName, appType string) (*sagemaker.DescribeAppOutput, error) {
	out, err := c.svc.DescribeApp(ctx, &sagemaker.DescribeAppInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfile),
		AppName:         aws.String(appName),
		AppType:         types.AppType(appType),
	})
	if err != nil {
		return nil, fmt.Errorf("sagemaker DescribeApp: %w", err)
	}
	return out, nil
}

// DeleteApp deletes a SageMaker app.
func (c *Client) DeleteApp(ctx context.Context, domainID, userProfile, appName, appType string) error {
	_, err := c.svc.DeleteApp(ctx, &sagemaker.DeleteAppInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfile),
		AppName:         aws.String(appName),
		AppType:         types.AppType(appType),
	})
	if err != nil {
		return fmt.Errorf("sagemaker DeleteApp: %w", err)
	}
	return nil
}

// CreatePresignedURL generates a presigned Studio URL for the given user.
func (c *Client) CreatePresignedURL(ctx context.Context, domainID, userProfile string) (string, error) {
	out, err := c.svc.CreatePresignedDomainUrl(ctx, &sagemaker.CreatePresignedDomainUrlInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfile),
	})
	if err != nil {
		return "", fmt.Errorf("sagemaker CreatePresignedDomainUrl: %w", err)
	}
	return aws.ToString(out.AuthorizedUrl), nil
}
