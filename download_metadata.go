package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func metadataDownloader(Metadataclient *MetadataBackend) {

	var (
		Bucket         = Metadataclient.Bucket // Download from this bucket
		Prefix         = "datasets/"           // Using this key prefix
		LocalDirectory = "web/data/"           // Into this directory
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
			err := downloadToFile(manager, LocalDirectory, Bucket, aws.ToString(obj.Key))
			if err != nil {
				log.Fatal("Error while downloading metadata files from metadata bucket", err)
			}
		}
		log.Infoln("Completed downloading metadata files")

	}

}

func downloadToFile(downloader *manager.Downloader, targetDirectory, bucket, key string) error {
	// Create the directories in the path
	file := filepath.Join(targetDirectory, key)
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
