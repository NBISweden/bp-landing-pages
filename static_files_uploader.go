package main

import (
	"context"
	"fmt"
	"path/filepath"

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

	result, err := uploadClient.UploadDirectory(context.TODO(), &transfermanager.UploadDirectoryInput{
		Bucket:    &bucket,
		KeyPrefix: new(""),
		Recursive: new(true),
		Source:    &localPath,
		Callback:  &typeFixer{},
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

type typeFixer struct{}

func (t *typeFixer) UpdateRequest(i *transfermanager.UploadObjectInput) {
	ext := filepath.Ext(*i.Key)

	switch ext {
	case ".css":
		i.ContentType = new("text/css")
	case ".svg":
		i.ContentType = new("image/svg+xml")
	case ".jpg", ".jpeg":
		i.ContentType = new("image/jpeg")
	case ".png":
		i.ContentType = new("image/png")
	case ".gif":
		i.ContentType = new("image/gif")
	case ".js":
		i.ContentType = new("application/javascript")
	case ".xml":
		i.ContentType = new("application/xml")
	case ".html":
		i.ContentType = new("text/html")

	}

	if i.ContentType != nil {
		log.Debugf("Content Type for %s: %s", *i.Key, *i.ContentType)
	} else {
		log.Debugf("no Content Type for %s", *i.Key)
	}
}
