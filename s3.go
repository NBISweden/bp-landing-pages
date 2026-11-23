package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type MetadataBackend struct {
	Client *s3.Client
	Bucket string
}
type DeploymentBackend struct {
	Client *s3.Client
	Bucket string
}

type Backend struct {
	Client *s3.Client
	Bucket string
}

func getBackend(s3config S3Config) *Backend {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3config.Region),
		config.WithHTTPClient(&http.Client{Transport: &http.Transport{}}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3config.AccessKey, s3config.SecretKey, "")),
	)

	if err != nil {
		log.Fatalf("Error while setting up s3 config: %v\n ", err)
	}
	client := s3.NewFromConfig(cfg,

		func(o *s3.Options) {
			o.BaseEndpoint = aws.String(s3config.URL)
			o.EndpointOptions.DisableHTTPS = strings.HasPrefix(s3config.URL, "http:")
			o.Region = s3config.Region
			o.UsePathStyle = true
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		})

	return &Backend{
		Client: client,
		Bucket: s3config.Bucket,
	}

}

func connectMetadatas3(mConf MetadataS3Config) *Backend {

	backend := getBackend(mConf.S3Config)
	resp, err := backend.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(backend.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the metadata bucket %v\n ", err)
	}
	log.Infoln("Connection established to metadata bucket", backend.Bucket)

	// Abort if 0 metadata files found in bucket
	numberOfFiles := 0
	for _, obj := range resp.Contents {
		if obj.Key != nil {
			numberOfFiles++
		}
	}
	if numberOfFiles == 0 {
		log.Fatal("No Metadata files found in bucket. Length of files ", numberOfFiles)
	}
	return backend
}

func connectDeployments3(dConf DeployS3Config) *Backend {

	backend := getBackend(dConf.S3Config)

	_, err := backend.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(backend.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the deplpyment bucket %v\n ", err)
	}
	log.Infoln("Connection established to deployment bucket", backend.Bucket)
	return backend
}
