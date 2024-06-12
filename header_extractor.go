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

// Function to extract header value from XML content
func getHeaderValueFromXMLContent(xmlContent []byte) (string, error) {
	var landingPageSet LandingPageSet
	err := xml.Unmarshal(xmlContent, &landingPageSet)
	if err != nil {
		log.Fatal("Unmarshalling XML file for header failed", err)
	}

	// Iterate over the string attributes to find the header value
	for _, attr := range landingPageSet.LandingPage.Attributes.StringAttributes {
		if attr.Tag == "header" {
			return attr.Value, nil
		}
	}
	return "sf", err

}

// Function to read the XML file and return its content
func readXMLFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
