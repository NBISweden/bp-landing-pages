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
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
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

func metadataDownloader(metadataclient *MetadataBackend) {

	var (
		Bucket         = metadataclient.Bucket // Download from this bucket
		Prefix         = "datasets/"           // Using this key prefix
		XMLDirectory   = "web/data/tmp/"       // Into this directory
		ImageDirectory = "web/static/img/"
	)
	client := Metadataclient.Client
	const maxImagesToDownload = 5
	stats := DownloadStats{}

	downloader := s3transfermanager.New(client, func(o *transfermanager.Options) {
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

func downloadToFile(manager *transfermanager.Client, targetDirectory, bucket, key string) error {
	// Create the directories in the path
	re := regexp.MustCompile("aa-Dataset-([a-z0-9]+)-([a-z0-9]+)")
	matches := re.FindStringSubmatch(key)
	if len(matches) == 0 {
		log.Infof("could not extract dataset ID from key: %s", key)
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

	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		log.Fatal("Failed to download metadatafiles", err)
	}
	return err
}

func downloadImageToFolder(downloader *transfermanager.Client, targetDirectory, bucket, key string) error {

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

	// Download the file
	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		log.Fatalf("failed to download image file: %v", err)
	}
	log.Debugf("Finished image %s", file)
	return nil
}
