package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func markDownCreator() {
	// Datasets directory path (change this to your XML directory path)
	DatasetDirectory := "web/data/datasets/"
	// Directory to save Markdown files
	markdownDir := "web/content/datasets/"

	// Create the Markdown directory if it doesn't exist
	if err := os.MkdirAll(markdownDir, 0755); err != nil {
		log.Fatal("Error creating directory:", err)
		os.Exit(1)
	}

	// Walk through the XML directory
	DatasetDirs, err := os.ReadDir(markdownDir)
	for _, e := range DatasetDirs {
		if e.IsDir() {
			filepath.WalkDir(filepath.Join(markdownDir, e.Name()), func(path string, d os.DirEntry, err error) error {
				if filepath.Ext(d.Name()) == ".xml" {
					fmt.Println(filepath.Rel(e.Name(), d.Name()))
					xmlpath := filepath.Join(markdownDir, e.Name(), "LANDING_PAGE", d.Name())
					fmt.Println("path ", xmlpath)
					xmlContent, err := os.Open(xmlpath)
					if err != nil {
						log.Error("error", err)
					}
					defer xmlContent.Close()
					if err != nil {
						log.Error("error", err)
					}
					outfile, err := os.Create(filepath.Join(DatasetDirectory, e.Name()+".xml"))
					if err != nil {
						log.Println("error writing xml file", err)
					}
					fmt.Println(xmlContent)
					_, err = io.Copy(outfile, xmlContent)
					if err != nil {
						log.Println("error writing xml file", err)
					}

				}

				return nil
			})
		}
	}
	if err != nil {
		log.Fatal("Error accessing file", err)
	}

	// Read XML file name
	log.Debug("getting info for file", DatasetDirectory)
	xmlContent, err := readXMLFile(DatasetDirectory)
	if err != nil {
		log.Fatalf("Error reading the XML file %V", err)
	}
	headerValue, doiValue, err := getHeaderValueFromXMLContent(xmlContent)
	if err != nil {
		log.Fatal("Error while getting header value from XML file", err)
	}

	log.Debugln("Header value: %V", headerValue)
	// Remove file extension
	fileNameWithoutExt := DatasetDirectory

	// Markdown content
	markdownContent := fmt.Sprintf(`---
title: "%s"
doi: "%s"
---

{{< datafetch variable="%s" >}}


Filename of the associated XML file: %s
`, headerValue, doiValue, fileNameWithoutExt, DatasetDirectory)

	// Create Markdown file
	mdFileName := filepath.Join(markdownDir, fileNameWithoutExt+".md")
	mdFile, err := os.Create(mdFileName)
	if err != nil {
		log.Fatal("Error creating Markdown file %V", err)
	}
	defer mdFile.Close()

	// Write Markdown content to the file
	_, err = mdFile.WriteString(markdownContent)
	if err != nil {
		log.Fatal("Error writing to Markdown file %V", mdFileName, err)
	}

	log.Debug("Markdown file %S created successfully!\n", mdFileName)

	if err != nil {
		log.Fatal("Error walking through directory:", err)
	}

}
