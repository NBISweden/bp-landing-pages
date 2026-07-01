package main

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
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
	pagefindCmd := exec.Command("pagefind", "--site", "web/public/", "--output-path", "web/static/pagefind/")
	pagefindCmd.Stdout = os.Stdout
	pagefindCmd.Stderr = os.Stderr
	pagefindErr := pagefindCmd.Run()
	if pagefindErr != nil {
		log.Fatal(pagefindErr)
	}
	log.Infof("Pagefind modules successfully built")
	dConf := getDeploymentConfig()
	DeploymenClient := connectDeployments3(dConf)
	staticSiteUploader(DeploymenClient)

}
