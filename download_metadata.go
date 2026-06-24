package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

const downloadTimeout = 60 * time.Second

type datasetIDExtractor struct {
	name    string
	pattern *regexp.Regexp
}

var datasetIDExtractors = []datasetIDExtractor{
	{
		name:    "standard aa-Dataset key",
		pattern: regexp.MustCompile(`(aa-Dataset-[a-z0-9]+-[a-z0-9]+)`),
	},
}

func isImageKey(key string) bool {
	ext := strings.ToLower(filepath.Ext(key))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func isXMLKey(key string) bool {
	return filepath.Ext(key) == ".xml"
}

func datasetIDFromKey(key string) (string, error) {
	for _, extractor := range datasetIDExtractors {
		if matches := extractor.pattern.FindStringSubmatch(key); len(matches) > 1 {
			log.Debugf("Matched dataset ID %q using extractor %q", matches[1], extractor.name)
			return matches[1], nil
		}
	}
	return "", errors.New("no dataset ID pattern matched key")
}

type DownloadStats struct {
	XMLDownloaded    int
	ImagesDownloaded int
	Skipped          int
	Failed           int
}

func metadataDownloader(Metadataclient *MetadataBackend) {

	var (
		Bucket         = Metadataclient.Bucket
		Prefix         = "datasets/"
		XmlDirectory   = "web/data/tmp/"
		ImageDirectory = "web/static/img/"
	)
	client := Metadataclient.Client
	const maxImagesToDownload = 5
	stats := DownloadStats{}

	downloader := s3manager.NewDownloader(client)
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
			if isXMLKey(key) {
				if err := downloadToFile(downloader, XmlDirectory, Bucket, key); err != nil {
					log.Warnf("Skipping XML %q: %v", key, err)
					stats.Failed++
				} else {
					stats.XMLDownloaded++
				}
				continue
			}

			// Handle image files (jpg, png, jpeg, gif, etc.)
			if isImageKey(key) {
				if stats.ImagesDownloaded >= maxImagesToDownload {
					stats.Skipped++
					continue
				}
				if err := downloadImageToFolder(downloader, ImageDirectory, Bucket, key); err != nil {
					log.Warnf("Skipping image %q: %v", key, err)
					stats.Failed++
				} else {
					stats.ImagesDownloaded++
				}
			}
		}
	}
	log.Infof("Completed: xml=%d images=%d skipped=%d failed=%d",
		stats.XMLDownloaded, stats.ImagesDownloaded, stats.Skipped, stats.Failed)
}

func downloadToFile(downloader *s3manager.Downloader, targetDirectory, bucket, key string) error {
	datasetID, err := datasetIDFromKey(key)
	if err != nil {
		return fmt.Errorf("skipping key %q: %w", key, err)
	}

	file := filepath.Join(targetDirectory, datasetID+".xml")
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return fmt.Errorf("creating directory for %q: %w", file, err)
	}

	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", file, err)
	}
	defer fd.Close()

	log.Debugf("Downloading XML %q -> %s", key, file)
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()
	if _, err = downloader.Download(ctx, fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key}); err != nil {
		return fmt.Errorf("downloading %q: %w", key, err)
	}
	log.Debugf("Finished XML %s", file)
	return nil
}

func downloadImageToFolder(downloader *s3manager.Downloader, targetDirectory, bucket, key string) error {
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid key structure: %q", key)
	}

	rootFolderName := parts[1]

	imageDir := filepath.Join(targetDirectory, rootFolderName)
	if err := os.MkdirAll(imageDir, 0775); err != nil {
		return err
	}

	// Get the original filename
	filename := filepath.Base(key)

	// Create the full file path
	file := filepath.Join(imageDir, filename)

	fd, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("creating image file %q: %w", file, err)
	}
	defer fd.Close()

	log.Debugf("Downloading image %q -> %s", key, file)
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()
	if _, err = downloader.Download(ctx, fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key}); err != nil {
		return fmt.Errorf("downloading image %q: %w", key, err)
	}
	log.Debugf("Finished image %s", file)
	return nil
}
