package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getConfig() MetadataS3Config {
	parseConfig()
	S3Conf := configS3Storage()
	return S3Conf

}

type MetadataS3Config struct {
	URL               string
	Port              int
	AccessKey         string
	SecretKey         string
	Bucket            string
	Region            string
	Chunksize         int
	Cacert            string
	WebMetadataFolder string
}

func configS3Storage() MetadataS3Config {
	s3 := MetadataS3Config{}
	s3.URL = viper.GetString("S3MetadataBucket.url")
	s3.AccessKey = viper.GetString("S3MetadataBucket.accesskey")
	s3.SecretKey = viper.GetString("S3MetadataBucket.secretkey")
	s3.Bucket = viper.GetString("S3MetadataBucket.bucket")
	s3.Port = 9000
	s3.WebMetadataFolder = viper.GetString("S3MetadataBucket.WebMetadataFolder")
	if viper.IsSet("s3.port") {
		s3.Port = viper.GetInt("S3MetadataBucket.port")
	}

	if viper.IsSet("s3.region") {
		s3.Region = viper.GetString("S3MetadataBucket.region")
	}

	if viper.IsSet("s3.chunksize") {
		s3.Chunksize = viper.GetInt("S3MetadataBucket.chunksize") * 1024 * 1024
	}

	if viper.IsSet("s3.cacert") {
		s3.Cacert = viper.GetString("S3MetadataBucket.cacert")
	}

	return s3
}

func parseConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	if viper.IsSet("configFile") {
		viper.SetConfigFile(viper.GetString("configFile"))
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Infoln("No config file found, using ENVs only")
		} else {
			log.Fatalf("Error when reading config file: '%s'", err)
		}
	}
}
