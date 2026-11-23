package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getMetadataConfig() MetadataS3Config {
	S3Conf := configS3Storage()
	return S3Conf
}

func getDeploymentConfig() DeployS3Config {
	dS3Conf := deployS3Storage()
	return dS3Conf
}

func setLogLevel() {
	switch viper.GetString("loglevel") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

type S3Config struct {
	URL       string
	Port      int
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Chunksize int
	Cacert    string
}

type MetadataS3Config struct {
	S3Config
	WebMetadataFolder string
}

type DeployS3Config struct {
	S3Config
}

func s3Config(prefix string) S3Config {
	s3 := S3Config{}
	s3.URL = viper.GetString(prefix + ".url")
	s3.AccessKey = viper.GetString(prefix + ".accesskey")
	s3.SecretKey = viper.GetString(prefix + ".secretkey")
	s3.Bucket = viper.GetString(prefix + ".bucket")

	viper.SetDefault(prefix+".port", 9000)
	viper.SetDefault(prefix+".region", "us-east-1")

	s3.Port = viper.GetInt(prefix + ".port")

	s3.Region = viper.GetString(prefix + ".region")

	s3.Chunksize = viper.GetInt("S3MetadataBucket.chunksize") * 1024 * 1024

	s3.Cacert = viper.GetString("S3MetadataBucket.cacert")

	return s3
}

func configS3Storage() MetadataS3Config {
	m := MetadataS3Config{
		s3Config("S3MetadataBucket"),
		viper.GetString("S3MetadataBucket.WebMetadataFolder"),
	}

	return m
}

func deployS3Storage() DeployS3Config {
	s3 := DeployS3Config{
		s3Config("S3DeploymentBucket"),
	}

	return s3
}

func prepareConfig() error {
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
			return fmt.Errorf("error when reading config file: %w", err)
		}
	}

	return nil
}
