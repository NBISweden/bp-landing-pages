package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type MetadataBackend struct {
	Client *s3.Client
	Bucket string
}

func connect_to_s3(mConf MetadataS3Config) *MetadataBackend {
	httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.MinVersion = tls.VersionTLS13
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(mConf.Region),
		config.WithHTTPClient(httpClient),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(mConf.AccessKey, mConf.SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: mConf.URL}, nil
			})),
	)

	if err != nil {
		log.Fatal("Could not connect to metadata S3 bucket: ", err)
	}
	client := s3.NewFromConfig(cfg)
	metadata_client := &MetadataBackend{

		Client: client,
		Bucket: mConf.Bucket,
	}

	return metadata_client
}
