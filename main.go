package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	parseConfig()
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

	log.Infoln("started app successfully")
	mConf := getMetadataConfig()
	Metadataclient := connectMetadatas3(mConf)
	log.Infof("Writing markdownfiles for respective XMLs")
	metadataDownloader(Metadataclient)
	markDownCreator()
	hugoCmd := exec.Command("hugo")
	hugoCmd.Dir = "./web/"
	hugoCmd.Stdout = os.Stdout
	hugoCmd.Stderr = os.Stderr
	err := hugoCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Hugo successfully built")

	dConf := getDeploymentConfig()
	DeploymenClient := connectDeployments3(dConf)
	staticSiteUploader(DeploymenClient)

}
