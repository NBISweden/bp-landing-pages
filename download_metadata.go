package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func metadataDownloader(Metadataclient *MetadataBackend) {

	var (
		Bucket         = Metadataclient.Bucket // Download from this bucket
		Prefix         = "datasets/"           // Using this key prefix
		XmlDirectory   = "web/data/tmp/"       // Into this directory
		ImageDirectory = "web/static/img/"
	)
	client := Metadataclient.Client

	manager := manager.NewDownloader(client)
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
				err := downloadToFile(manager, XmlDirectory, Bucket, key)
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

func downloadToFile(downloader *manager.Downloader, targetDirectory, bucket, key string) error {
	// Create the directories in the path
	re := regexp.MustCompile("aa-Dataset-([a-z0-9]+)-([a-z0-9]+)")
	matches := re.FindStringSubmatch(key)
	if len(matches) == 0 {
		log.Infof("could not extract dataset ID from key: %s", key)
	}
	DatasetId := matches[0]
	fileStr := fmt.Sprintf("%s%s", DatasetId, ".xml")
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

	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		log.Fatal("Failed to download metadatafiles", err)
	}
	return err
}

func downloadImageToFolder(downloader *manager.Downloader, targetDirectory, bucket, key string) error {

	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		log.Fatalf("invalid key structure: %s", key)
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
		log.Fatalf("error while creating image file: %w", err)
	}
	defer fd.Close()

	// Download the file
	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		log.Fatalf("failed to download image file: %w", err)
	}

	log.Infof("Downloaded image to: %s", file)
	return nil
}
