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
		LocalDirectory = "web/content/"        // Into this directory
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
			log.Fatalln("error:", err)
		}
		for _, obj := range page.Contents {
			if err := downloadToFile(manager, LocalDirectory, Bucket, aws.ToString(obj.Key)); err != nil {
				log.Fatalln("error:", err)
			}
		}

	}
	log.Infoln("Completed downloading metadatafiles")

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
		return err
	}
	defer fd.Close()
	// Download the file using the AWS SDK for Go
	//fmt.Printf("Downloading s3://%s/%s to %s...\n", bucket, key, file)

	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		log.Infoln("Failed to download metadatafiles")
	}
	return err
}
