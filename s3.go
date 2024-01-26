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

type metadata_backend struct {
	Client *s3.Client
	Bucket string
}

func connect_to_s3(mconf metadata_s3_config) *metadata_backend {
	httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.MinVersion = tls.VersionTLS13
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(mconf.Region),
		config.WithHTTPClient(httpClient),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(mconf.AccessKey, mconf.SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: mconf.URL}, nil
			})),
	)

	if err != nil {
		log.Fatal("Could not connect to metadata S3 bucket: ", err)
	}
	client := s3.NewFromConfig(cfg)
	metadata_client := &metadata_backend{

		Client: client,
		Bucket: mconf.Bucket,
	}

	return metadata_client
}
