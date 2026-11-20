package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
)

func staticSiteUploader(deploymenClient *DeploymentBackend) error {
	var (
		localPath = "web/public/"
		bucket    = deploymenClient.Bucket
	)

	client := deploymenClient.Client
	uploadClient := transfermanager.New(client)

	recurse := true
	keyPrefix := ""

	result, err := uploadClient.UploadDirectory(context.TODO(), &transfermanager.UploadDirectoryInput{
		Bucket:    &bucket,
		KeyPrefix: &keyPrefix,
		Recursive: &recurse,
		Source:    &localPath,
	},
	)

	if err != nil {
		return fmt.Errorf("failed to upload %s: %w", localPath, err)
	}

	if result.ObjectsFailed != 0 {
		return fmt.Errorf("%d objects failed to upload", result.ObjectsFailed)
	}

	log.Infof("Successfully uploaded built static site to the bucket (%d objects)", result.ObjectsUploaded)
	return nil
}
