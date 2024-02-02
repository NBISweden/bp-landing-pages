package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("started app successfully")
	mConf := getConfig()
	Metadataclient := connect_to_s3(mConf)
	log.Infof("Connection to the bucket established")
	metadataDownloader(Metadataclient)
}
