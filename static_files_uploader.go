package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func test(DeploymenClient *DeploymentBackend) {
	var (
		localPath = "web/public/"
		bucket    = DeploymenClient.Bucket
	)

	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localPath, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	// For each file found walking, upload it to Amazon S3
	for path := range walker {
		rel, err := filepath.Rel(localPath, path)
		if err != nil {
			log.Fatalln("Unable to get relative path:", path, err)
		}
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}
		buf := make([]byte, 512)
		_, err = file.Read(buf)
		contentType := http.DetectContentType(buf)

		if err != nil {
			log.Fatalln("File bytes empty", path, err)
		}

		client := DeploymenClient.Client
		uploader := manager.NewUploader(client)
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:      &bucket,
			Key:         aws.String(rel),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		println(contentType)
		file.Close()
		if err != nil {
			log.Fatalln("Failed to upload", path, err)
		}
		log.Infoln("Uploaded", path, result.Location)
	}
}

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}
