package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func metadataDownloader(metadataclient *MetadataBackend) {

	var (
		Bucket         = metadataclient.Bucket // Download from this bucket
		Prefix         = "datasets/"           // Using this key prefix
		XMLDirectory   = "web/data/tmp/"       // Into this directory
		ImageDirectory = "web/static/img/"
	)
	client := metadataclient.Client

	manager := transfermanager.New(client, func(o *transfermanager.Options) {
		o.PartSizeBytes = 64 * 1024 * 1024
		o.Concurrency = 3
	})
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &Bucket,
		Prefix: &Prefix,
	})
	log.Infoln("Started downloading metadata files")
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalln("Error while paginating the s3 bucket", err)
		}
		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)

			// Handle XML files
			if filepath.Ext(key) == ".xml" {
				err := downloadToFile(manager, XMLDirectory, Bucket, key)
				if err != nil {
					log.Fatal("Error while downloading metadata files from metadata bucket", err)
				}
			}
			// Handle image files (jpg, png, jpeg, gif, etc.)
			ext := strings.ToLower(filepath.Ext(key))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				err := downloadImageToFolder(manager, ImageDirectory, Bucket, key)
				if err != nil {
					log.Fatal("Error while downloading image files from metadata bucket", err)
				}

			}
		}
	}
	log.Infoln("Completed downloading metadata files")
}

func downloadToFile(manager *transfermanager.Client, targetDirectory, bucket, key string) error {
	// Create the directories in the path
	re := regexp.MustCompile("aa-Dataset-([a-z0-9]+)-([a-z0-9]+)")
	matches := re.FindStringSubmatch(key)
	if len(matches) == 0 {
		log.Infof("could not extract dataset ID from key: %s", key)
	}
	DatasetID := matches[0]
	fileStr := fmt.Sprintf("%s%s", DatasetID, ".xml")
	file := filepath.Join(targetDirectory, fileStr)
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		log.Fatal("Error while writing XML files to folder", err)
	}
	defer fd.Close()

	_, err = manager.DownloadObject(context.TODO(), &transfermanager.DownloadObjectInput{Bucket: &bucket, Key: &key})
	return err
}

func downloadImageToFolder(downloader *transfermanager.Client, targetDirectory, bucket, key string) error {

	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		log.Errorf("invalid key structure: %s", key)
	}

	rootFolderName := parts[1] // This gets the folder name after "datasets/"

	// Create the target directory: targetDirectory/rootFolderName/
	imageDir := filepath.Join(targetDirectory, rootFolderName)
	if err := os.MkdirAll(imageDir, 0775); err != nil {
		return err
	}

	// Get the original filename
	filename := filepath.Base(key)

	// Create the full file path
	file := filepath.Join(imageDir, filename)

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		log.Errorf("error while creating image file: %v", err)
	}
	defer fd.Close()

	// Download the file
	_, err = downloader.DownloadObject(context.TODO(), &transfermanager.DownloadObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return err
	}

	log.Infof("Downloaded image to: %s", file)
	return nil
}
