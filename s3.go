package main

import (
	"context"
	"crypto/tls"
	"net/http"

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
		config.WithRegion("auto"),
		config.WithHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(mConf.AccessKey, mConf.SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: mConf.URL}, nil
			})),
	)

	if err != nil {
		log.Fatalf("Error while setting up s3 config: %v\n ", err)
	}
	client := s3.NewFromConfig(cfg)
	metadata_client := &MetadataBackend{

		Client: client,
		Bucket: mConf.Bucket,
	}
	_, err = metadata_client.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(metadata_client.Bucket),
	})
	if err != nil {
		log.Fatalf("Error while connecting to the metadata bucket %v\n ", err)
	} else {
		log.Infoln("Connection established to metadata bucket", metadata_client.Bucket)
	}
	return metadata_client
}

func connectDeployments3(dConf DeployS3Config) *DeploymentBackend {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(dConf.AccessKey, dConf.SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: dConf.URL}, nil
			})),
	)

	if err != nil {
		log.Fatalf("Error while setting up s3 config: %v\n ", err)
	}
	client := s3.NewFromConfig(cfg)
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
