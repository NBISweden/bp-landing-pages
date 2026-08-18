package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type LandingPageSet struct {
	Pages []LandingPage `xml:"LANDING_PAGE"`
}

type LandingPage struct {
	Alias            string           `xml:"alias,attr"`
	DatasetRef       DatasetRef       `xml:"DATASET_REF"`
	SampleImageFiles SampleImageFiles `xml:"SAMPLE_IMAGE_FILES"`
	Attributes       Attributes       `xml:"ATTRIBUTES"`
	RemsLink         string           `xml:"REMS_ACCESS_LINK"`
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

func normalizeSampleImageName(name string) string {
	name = strings.TrimPrefix(name, "LANDING_PAGE/THUMBNAILS/")
	name = strings.TrimPrefix(name, "LANDING_PAGE/")
	return strings.TrimSuffix(name, ".c4gh")
}

func frontMatterFromAttributes(lp LandingPage, fileNameWithoutExt string) map[string]any {
	attrs := lp.Attributes
	result := map[string]any{}

	stringFields := []string{
		"header",
		"doi",
		"dataset_title",
		"dataset_short_name",
		"dataset_version",
		"metadata_standard_version",
		"center_name",
		"access_approval_process",
		"type_of_access",
		"allowed_geographical_distribution",
		"duration_of_use",
		"defined_research_question_required",
		"informed_consent_form_defined_use_restrictions",
		"custom_use_restrictions",
		"policy_text",
		"dataset_description",
	}
	for _, key := range stringFields {
		result[key] = attrs.GetString(key)
	}
	result["rems_access_link"] = lp.RemsLink

	numericFields := []string{
		"number_of_biological_beings",
		"number_of_cases",
		"number_of_wsis",
		"number_of_observations",
		"number_of_annotations",
		"dataset_size",
		"year_of_submission",
	}
	for _, key := range numericFields {
		result[key] = attrs.GetNumber(key)
	}

	setFields := []string{
		"keywords",
		"animal_species",
		"anatomical_sites",
		"age_at_extractions",
		"extraction_methods",
		"specimen_types",
		"stainings",
		"medical_diagnoses",
		"image_types",
		"image_resolutions",
		"geographical_areas",
		"changelog",
		"cite_as",
		"references",
		"comments",
		"allowed_uses",
	}
	for _, key := range setFields {
		values := attrs.GetSet(key)
		if len(values) > 0 {
			result[key] = values
		}
	}

	if len(lp.SampleImageFiles.Files) > 0 {
		var sampleImages []map[string]any
		for _, f := range lp.SampleImageFiles.Files {
			cleanName := normalizeSampleImageName(f.Filename)
			entry := map[string]any{
				"filename": path.Join("/img", fileNameWithoutExt, cleanName),
				"filetype": f.Filetype,
			}
			if f.Checksum != "" {
				entry["checksum"] = f.Checksum
			}
			if f.UnencryptedChecksum != "" {
				entry["unencrypted_checksum"] = f.UnencryptedChecksum
			}
			sampleImages = append(sampleImages, entry)
		}
		result["sample_images"] = sampleImages
	}

	return result
}

func toFrontMatter(lp LandingPage, fileNameWithoutExt string) (string, error) {
	data := frontMatterFromAttributes(lp, fileNameWithoutExt)
	out, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal front matter: %w", err)
	}
	return "---\n" + string(out) + "---\n", nil
}

func markdownWriter(xmlContent []byte, fileNameWithoutExt string) (string, error) {
	var set LandingPageSet
	if err := xml.Unmarshal(xmlContent, &set); err != nil {
		return "", fmt.Errorf("unmarshal landing page XML: %w", err)
	}
	if len(set.Pages) == 0 {
		return "", fmt.Errorf("no landing pages found in XML")
	}

	front, err := toFrontMatter(set.Pages[0], fileNameWithoutExt)
	if err != nil {
		return "", err
	}

	return front, nil
}
