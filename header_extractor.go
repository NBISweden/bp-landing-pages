package main

import (
	"encoding/xml"
	"log"
	"os"
)

// Struct definitions
type LandingPageSet struct {
	LandingPage LandingPage `xml:"LANDING_PAGE"`
}

type LandingPage struct {
	Attributes Attributes `xml:"ATTRIBUTES"`
}

type Attributes struct {
	StringAttributes []StringAttribute `xml:"STRING_ATTRIBUTE"`
}

type StringAttribute struct {
	Tag   string `xml:"TAG"`
	Value string `xml:"VALUE"`
}

func getHeaderValueFromXMLContent(xmlContent []byte) (string, string, error) {
	var landingPageSet LandingPageSet
	err := xml.Unmarshal(xmlContent, &landingPageSet)
	if err != nil {
		log.Fatal("Unmarshalling XML file failed", err)
	}

	var header, doi string

	// Iterate over the string attributes to find the header and doi values
	for _, attr := range landingPageSet.LandingPage.Attributes.StringAttributes {
		switch attr.Tag {
		case "header":
			header = attr.Value
		case "doi":
			doi = attr.Value
		}
	}

	if header == "" {
		return "", "", err
	}

	return header, doi, nil
}

// Function to read the XML file and return its content
func readXMLFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
