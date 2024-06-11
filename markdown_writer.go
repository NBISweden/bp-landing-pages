package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func markDownCreator() {
	// XML directory path (change this to your XML directory path)
	xmlDirPath := "web/data/datasets/"

	// Directory to save Markdown files
	markdownDir := "web/content/datasets/"

	// Create the Markdown directory if it doesn't exist
	if err := os.MkdirAll(markdownDir, 0755); err != nil {
		log.Fatal("Error creating directory:", err)
		os.Exit(1)
	}

	// Walk through the XML directory
	err := filepath.Walk(xmlDirPath, func(xmlFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("Error accessing file", err)
			return nil
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read XML file name
		xmlFileName := filepath.Base(xmlFilePath)
		fmt.Println(xmlFilePath)
		xmlContent, err := readXMLFile(xmlFilePath)
		if err != nil {
			fmt.Println(err)
		}
		headerValue, err := getHeaderValueFromXMLContent(xmlContent)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Header value: %s\n", headerValue)
		// Remove file extension
		fileNameWithoutExt := strings.TrimSuffix(xmlFileName, filepath.Ext(xmlFileName))

		// Markdown content
		markdownContent := fmt.Sprintf(`---
title: "%s"
---

{{< datafetch variable="%s" >}}


Filename of the associated XML file: %s
`, headerValue, fileNameWithoutExt, xmlFileName)

		// Create Markdown file
		mdFileName := filepath.Join(markdownDir, fileNameWithoutExt+".md")
		mdFile, err := os.Create(mdFileName)
		if err != nil {
			log.Fatal("Error creating Markdown file %V", err)
			return nil
		}
		defer mdFile.Close()

		// Write Markdown content to the file
		_, err = mdFile.WriteString(markdownContent)
		if err != nil {
			log.Fatal("Error writing to Markdown file %V", mdFileName, err)
			return nil
		}

		log.Infoln("Markdown file %S created successfully!\n", mdFileName)

		return nil
	})

	if err != nil {
		log.Fatal("Error walking through directory:", err)
	}
}
