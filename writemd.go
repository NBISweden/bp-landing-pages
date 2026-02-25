package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type LandingPageSet struct {
	Pages []LandingPage `xml:"LANDING_PAGE"`
}

type LandingPage struct {
	Alias            string           `xml:"alias,attr"`
	DatasetRef       DatasetRef       `xml:"DATASET_REF"`
	SampleImageFiles SampleImageFiles `xml:"SAMPLE_IMAGE_FILES"`
	Attributes       Attributes       `xml:"ATTRIBUTES"`
}

type DatasetRef struct {
	Alias string `xml:"alias,attr"`
}

type SampleImageFiles struct {
	Files []SampleImageFile `xml:"SAMPLE_IMAGE_FILE"`
}

type SampleImageFile struct {
	Filename            string `xml:"filename,attr"`
	Filetype            string `xml:"filetype,attr"`
	ChecksumMethod      string `xml:"checksum_method,attr"`
	Checksum            string `xml:"checksum,attr"`
	UnencryptedChecksum string `xml:"unencrypted_checksum,attr"`
}

type Attributes struct {
	Strings []StringAttr  `xml:"STRING_ATTRIBUTE"`
	Numbers []NumericAttr `xml:"NUMERIC_ATTRIBUTE"`
	Sets    []SetAttr     `xml:"SET_ATTRIBUTE"`
}

type StringAttr struct {
	Tag   string `xml:"TAG"`
	Value string `xml:"VALUE"`
}

type NumericAttr struct {
	Tag   string `xml:"TAG"`
	Value string `xml:"VALUE"`
}

type SetAttr struct {
	Tag   string   `xml:"TAG"`
	Value SetValue `xml:"VALUE"`
}

type SetValue struct {
	Items []StringAttr `xml:"STRING_ATTRIBUTE"`
}

func (a Attributes) GetString(tag string) string {
	for _, s := range a.Strings {
		if s.Tag == tag {
			return s.Value
		}
	}
	return ""
}

func (a Attributes) GetNumber(tag string) string {
	for _, n := range a.Numbers {
		if n.Tag == tag {
			return n.Value
		}
	}
	return ""
}

func (a Attributes) GetSet(tag string) []string {
	for _, s := range a.Sets {
		if s.Tag == tag {
			var list []string
			for _, item := range s.Value.Items {
				if strings.TrimSpace(item.Value) != "" {
					list = append(list, item.Value)
				}
			}
			return list
		}
	}
	return nil
}

func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "\t", "    ") // Replace tabs with 4 spaces
	s = strings.ReplaceAll(s, "\"", "'")
	return s
}
func writeStringField(b *strings.Builder, key, value string) {
	if strings.TrimSpace(value) == "" {
		fmt.Fprintf(b, "%s: \"\"\n", key)
	} else {
		fmt.Fprintf(b, "%s: \"%s\"\n", key, escapeYAML(value))
	}
}

func writeListField(b *strings.Builder, key string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(b, "%s:\n", key)
	for _, item := range values {
		fmt.Fprintf(b, "  - \"%s\"\n", escapeYAML(item))
	}
}

func toFrontMatter(lp LandingPage, fileNameWithoutExt string) string {
	attrs := lp.Attributes
	var b strings.Builder
	b.WriteString("---\n")

	// Extract fields individually in desired order
	writeStringField(&b, "header", attrs.GetString("header"))
	writeStringField(&b, "doi", attrs.GetString("doi"))
	writeStringField(&b, "dataset_title", attrs.GetString("dataset_title"))
	writeStringField(&b, "dataset_short_name", attrs.GetString("dataset_short_name"))
	writeStringField(&b, "dataset_version", attrs.GetString("dataset_version"))
	writeStringField(&b, "metadata_standard_version", attrs.GetString("metadata_standard_version"))
	writeStringField(&b, "center_name", attrs.GetString("center_name"))
	writeStringField(&b, "access_approval_process", attrs.GetString("access_approval_process"))
	writeStringField(&b, "type_of_access", attrs.GetString("type_of_access"))
	writeStringField(&b, "allowed_geographical_distribution", attrs.GetString("allowed_geographical_distribution"))
	writeStringField(&b, "duration_of_use", attrs.GetString("duration_of_use"))
	writeStringField(&b, "defined_research_question_required", attrs.GetString("defined_research_question_required"))
	writeStringField(&b, "informed_consent_form_defined_use_restrictions", attrs.GetString("informed_consent_form_defined_use_restrictions"))
	writeStringField(&b, "custom_use_restrictions", attrs.GetString("custom_use_restrictions"))
	writeStringField(&b, "policy_text", attrs.GetString("policy_text"))

	// Numeric fields
	writeStringField(&b, "number_of_biological_beings", attrs.GetNumber("number_of_biological_beings"))
	writeStringField(&b, "number_of_cases", attrs.GetNumber("number_of_cases"))
	writeStringField(&b, "number_of_wsis", attrs.GetNumber("number_of_wsis"))
	writeStringField(&b, "number_of_observations", attrs.GetNumber("number_of_observations"))
	writeStringField(&b, "number_of_annotations", attrs.GetNumber("number_of_annotations"))
	writeStringField(&b, "dataset_size", attrs.GetNumber("dataset_size"))
	writeStringField(&b, "year_of_submission", attrs.GetNumber("year_of_submission"))
	writeStringField(&b, "dataset_description", attrs.GetString("dataset_description"))
	// Set fields
	writeListField(&b, "keywords", attrs.GetSet("keywords"))
	writeListField(&b, "animal_species", attrs.GetSet("animal_species"))
	writeListField(&b, "anatomical_sites", attrs.GetSet("anatomical_sites"))
	writeListField(&b, "age_at_extractions", attrs.GetSet("age_at_extractions"))
	writeListField(&b, "extraction_methods", attrs.GetSet("extraction_methods"))
	writeListField(&b, "specimen_types", attrs.GetSet("specimen_types"))
	writeListField(&b, "stainings", attrs.GetSet("stainings"))
	writeListField(&b, "medical_diagnoses", attrs.GetSet("medical_diagnoses"))
	writeListField(&b, "image_types", attrs.GetSet("image_types"))
	writeListField(&b, "image_resolutions", attrs.GetSet("image_resolutions"))
	writeListField(&b, "geographical_areas", attrs.GetSet("geographical_areas"))
	writeListField(&b, "changelog", attrs.GetSet("changelog"))
	writeListField(&b, "cite_as", attrs.GetSet("cite_as"))
	writeListField(&b, "references", attrs.GetSet("references"))
	writeListField(&b, "comments", attrs.GetSet("comments"))
	writeListField(&b, "allowed_uses", attrs.GetSet("allowed_uses"))

	// Sample images
	if len(lp.SampleImageFiles.Files) > 0 {
		b.WriteString("sample_images:\n")
		for _, f := range lp.SampleImageFiles.Files {
			cleanName := strings.TrimPrefix(f.Filename, "LANDING_PAGE/THUMBNAILS/")
			cleanName = strings.TrimPrefix(cleanName, "LANDING_PAGE/")
			cleanName = strings.TrimSuffix(cleanName, ".c4gh")

			finalPath := path.Join("/img", fileNameWithoutExt, cleanName)

			b.WriteString(fmt.Sprintf("  - filename: %q\n", finalPath))

			b.WriteString("    filetype: \"" + f.Filetype + "\"\n")
			if f.Checksum != "" {
				b.WriteString("    checksum: \"" + f.Checksum + "\"\n")
			}
			if f.UnencryptedChecksum != "" {
				b.WriteString("    unencrypted_checksum: \"" + f.UnencryptedChecksum + "\"\n")
			}
		}
	}

	b.WriteString("---\n\n")
	return b.String()
}

func markdownWriter(xmlContent []byte, fileNameWithoutExt string) (string, error) {

	var set LandingPageSet
	if err := xml.Unmarshal(xmlContent, &set); err != nil {
		log.Errorf("Error while unmarshallaing %d", err)
	}

	front := toFrontMatter(set.Pages[0], fileNameWithoutExt)

	return front, nil
}
