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
	metadataClient := &MetadataBackend{

		Client: client,
		Bucket: mConf.Bucket,
	}
	resp, err := metadataClient.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(metadataClient.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the metadata bucket %v\n ", err)
	}
	log.Infoln("Connection established to metadata bucket", metadataClient.Bucket)

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
	return metadataClient
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
	deploymentClient := &DeploymentBackend{

		Client: client,
		Bucket: dConf.Bucket,
	}
	_, err = deploymentClient.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(deploymentClient.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the deplpyment bucket %v\n ", err)
	}
	log.Infoln("Connection established to deployment bucket", deploymentClient.Bucket)
	return deploymentClient
}
