package main

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("started app successfully")
	mConf := getMetadataConfig()
	Metadataclient := connectMetadatas3(mConf)
	log.Infof("Connection to the bucket established")
	metadataDownloader(Metadataclient)
	markDownCreator()
	cmd := exec.Command("hugo")
	cmd.Dir = "./web/"
	cmd.Run()
	log.Infof("Hugo successfully built")
	dConf := getDeploymentConfig()
	DeploymenClient := connectDeployments3(dConf)
	log.Infof("Connection to the bucket established")
	test(DeploymenClient)

}
