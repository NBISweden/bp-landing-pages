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
	hugo_cmd := exec.Command("hugo")
	hugo_cmd.Dir = "./web/"
	hugo_cmd.Stdout = os.Stdout
	hugo_cmd.Stderr = os.Stderr
	err := hugo_cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Hugo successfully built")
	pagefind_cmd := exec.Command("pagefind", "--site", "web/public/", "--output-path", "web/static/pagefind/")
	pagefind_cmd.Stdout = os.Stdout
	pagefind_cmd.Stderr = os.Stderr
	pagefind_err := pagefind_cmd.Run()
	if pagefind_err != nil {
		log.Fatal(pagefind_err)
	}
	log.Infof("Pagefind modules successfully built")
	dConf := getDeploymentConfig()
	DeploymenClient := connectDeployments3(dConf)
	staticSiteUploader(DeploymenClient)

}
