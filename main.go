package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("started app successfully")
	mconf := get_config()
	m_client := connect_to_s3(mconf)
	log.Infof("Connection to the bucket established")
	metadata_downloader(m_client)
}
