package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func markDownCreator() {
	// XML directory path (change this to your XML directory path)
	xmlDirPath := "web/data/datasets/"

	// Directory to save Markdown files
	markdownDir := "web/content/datasets/"

	// Create the Markdown directory if it doesn't exist
	if err := os.MkdirAll(markdownDir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	// Walk through the XML directory
	err := filepath.Walk(xmlDirPath, func(xmlFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", xmlFilePath, err)
			return nil
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read XML file name
		xmlFileName := filepath.Base(xmlFilePath)

		// Remove file extension
		fileNameWithoutExt := strings.TrimSuffix(xmlFileName, filepath.Ext(xmlFileName))

		// Markdown content
		markdownContent := fmt.Sprintf(`
---
---

{{< datafetch variable="%s" >}}


Filename of the associated XML file: %s
`, fileNameWithoutExt, xmlFileName)

		// Create Markdown file
		mdFileName := filepath.Join(markdownDir, fileNameWithoutExt+".md")
		mdFile, err := os.Create(mdFileName)
		if err != nil {
			fmt.Printf("Error creating Markdown file %s: %v\n", mdFileName, err)
			return nil
		}
		defer mdFile.Close()

		// Write Markdown content to the file
		_, err = mdFile.WriteString(markdownContent)
		if err != nil {
			fmt.Printf("Error writing to Markdown file %s: %v\n", mdFileName, err)
			return nil
		}

		fmt.Printf("Markdown file %s created successfully!\n", mdFileName)

		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
		os.Exit(1)
	}
}
