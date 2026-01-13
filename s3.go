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

func connectMetadatas3(mConf MetadataS3Config) *MetadataBackend {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithHTTPClient(&http.Client{Transport: &http.Transport{}}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(mConf.AccessKey, mConf.SecretKey, "")),
	)

	if err != nil {
		log.Fatalf("Error while setting up s3 config: %v\n ", err)
	}
	client := s3.NewFromConfig(cfg,

		func(o *s3.Options) {
			o.BaseEndpoint = aws.String(mConf.URL)
			o.EndpointOptions.DisableHTTPS = strings.HasPrefix(mConf.URL, "http:")
			o.Region = "us-west-1"
			o.UsePathStyle = true
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		})
	metadata_client := &MetadataBackend{

		Client: client,
		Bucket: mConf.Bucket,
	}
	resp, err := metadata_client.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(metadata_client.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the metadata bucket %v\n ", err)
	} else {
		log.Infoln("Connection established to metadata bucket", metadata_client.Bucket)
	}

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
	return metadata_client
}

func connectDeployments3(dConf DeployS3Config) *DeploymentBackend {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(dConf.AccessKey, dConf.SecretKey, "")),
	)

	if err != nil {
		log.Fatalf("Error while setting up s3 config: %v\n ", err)
	}
	client := s3.NewFromConfig(cfg,
		func(o *s3.Options) {
			o.BaseEndpoint = aws.String(dConf.URL)
			o.EndpointOptions.DisableHTTPS = strings.HasPrefix(dConf.URL, "http:")
			o.Region = "us-east-1"
			o.UsePathStyle = true
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		})
	deployment_client := &DeploymentBackend{

		Client: client,
		Bucket: dConf.Bucket,
	}
	_, err = deployment_client.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(deployment_client.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the deplpyment bucket %v\n ", err)
	} else {
		log.Infoln("Connection established to deployment bucket", deployment_client.Bucket)
	}
	return deployment_client
}
