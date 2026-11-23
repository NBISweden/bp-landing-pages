package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func main() {

	err := prepareConfig()
	if err != nil {
		log.Fatalf("failed to initiate configuration: %v", err)
	}

	setLogLevel()

	log.Infoln("started app successfully")
	mConf := getMetadataConfig()
	Metadataclient := connectMetadatas3(mConf)
	log.Infof("Writing markdownfiles for respective XMLs")

	err = metadataDownloader(Metadataclient)
	if err != nil {
		log.Fatalf("metadata download failed: %v", err)
	}

	markDownCreator()
	hugoCmd := exec.Command("hugo")
	hugoCmd.Dir = "./web/"
	hugoCmd.Stdout = os.Stdout
	hugoCmd.Stderr = os.Stderr
	err = hugoCmd.Run()
	if err != nil {
		log.Fatalf("hugo run failed: %v", err)
	}
	log.Infof("Hugo successfully built")

	dConf := getDeploymentConfig()
	DeploymenClient := connectDeployments3(dConf)

	err = staticSiteUploader(DeploymenClient)

	if err != nil {
		log.Fatalf("upload failed: %v", err)
	}
}
