package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func get_config() metadata_s3_config {
	parseConfig()
	s3_conf := configS3Storage()
	return s3_conf

}

type metadata_s3_config struct {
	URL                 string
	Port                int
	AccessKey           string
	SecretKey           string
	Bucket              string
	Region              string
	Chunksize           int
	Cacert              string
	Web_metadata_folder string
}

func configS3Storage() metadata_s3_config {
	s3 := metadata_s3_config{}
	s3.URL = viper.GetString("metadata_s3.url")
	s3.AccessKey = viper.GetString("metadata_s3.accesskey")
	s3.SecretKey = viper.GetString("metadata_s3.secretkey")
	s3.Bucket = viper.GetString("metadata_s3.bucket")
	s3.Port = 9000
	s3.Web_metadata_folder = viper.GetString("metadata_s3.web_metadata_folder")
	if viper.IsSet("s3.port") {
		s3.Port = viper.GetInt("metadata_s3.port")
	}

	if viper.IsSet("s3.region") {
		s3.Region = viper.GetString("metadata_s3.region")
	}

	if viper.IsSet("s3.chunksize") {
		s3.Chunksize = viper.GetInt("metadata_s3.chunksize") * 1024 * 1024
	}

	if viper.IsSet("s3.cacert") {
		s3.Cacert = viper.GetString("metadata_s3.cacert")
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
