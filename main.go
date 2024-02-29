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
	cmd := exec.Command("hugo")
	cmd.Dir = "./web/"
	cmd.Run()
	log.Infof("Hugo successfully built")
	deConf := getDeploymentConfig()
	connectDeployments3(deConf)

}
