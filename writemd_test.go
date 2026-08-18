package main

import (
	"strings"
	"testing"
)

func TestMarkdownWriterRejectsEmptyLandingPages(t *testing.T) {
	_, err := markdownWriter([]byte(`<LANDING_PAGE_SET></LANDING_PAGE_SET>`), "dataset")
	if err == nil {
		t.Fatal("expected an error for an empty landing page set")
	}
}

func TestMarkdownWriterProducesFrontMatterFromSyntheticXML(t *testing.T) {
	xmlData := []byte(`
		<?xml version="1.0" encoding="utf-8"?>
		<LANDING_PAGE_SET>
			<LANDING_PAGE alias="LANDING_PAGE_sample">
				<DATASET_REF alias="DATASET_sample"></DATASET_REF>
				<SAMPLE_IMAGE_FILES>
					<SAMPLE_IMAGE_FILE filename="LANDING_PAGE/THUMBNAILS/IMAGE_sample.jpg" filetype="jpg" checksum_method="SHA256" checksum="abc123" unencrypted_checksum="def456"></SAMPLE_IMAGE_FILE>
				</SAMPLE_IMAGE_FILES>
				<ATTRIBUTES>
					<STRING_ATTRIBUTE>
						<TAG>header</TAG>
						<VALUE>Sample cohort header</VALUE>
					</STRING_ATTRIBUTE>
					<STRING_ATTRIBUTE>
						<TAG>doi</TAG>
						<VALUE xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:nil="true"></VALUE>
					</STRING_ATTRIBUTE>
					<STRING_ATTRIBUTE>
						<TAG>dataset_title</TAG>
						<VALUE>Example synthetic dataset</VALUE>
					</STRING_ATTRIBUTE>
					<STRING_ATTRIBUTE>
						<TAG>dataset_description</TAG>
						<VALUE>A short description with a newline.
This is the second line.</VALUE>
					</STRING_ATTRIBUTE>
					<STRING_ATTRIBUTE>
						<TAG>policy_text</TAG>
						<VALUE>Example policy text:
This is a multiline policy summary.</VALUE>
					</STRING_ATTRIBUTE>
					<NUMERIC_ATTRIBUTE>
						<TAG>number_of_cases</TAG>
						<VALUE>42</VALUE>
					</NUMERIC_ATTRIBUTE>
					<NUMERIC_ATTRIBUTE>
						<TAG>dataset_size</TAG>
						<VALUE xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:nil="true"></VALUE>
					</NUMERIC_ATTRIBUTE>
					<SET_ATTRIBUTE>
						<TAG>keywords</TAG>
						<VALUE>
							<STRING_ATTRIBUTE>
								<TAG>keyword</TAG>
								<VALUE>sample</VALUE>
							</STRING_ATTRIBUTE>
							<STRING_ATTRIBUTE>
								<TAG>keyword</TAG>
								<VALUE>synthetic</VALUE>
							</STRING_ATTRIBUTE>
						</VALUE>
					</SET_ATTRIBUTE>
					<SET_ATTRIBUTE>
						<TAG>allowed_uses</TAG>
						<VALUE>
							<STRING_ATTRIBUTE>
								<TAG>allowed_use</TAG>
								<VALUE>Research use only</VALUE>
							</STRING_ATTRIBUTE>
						</VALUE>
					</SET_ATTRIBUTE>
				</ATTRIBUTES>
				<REMS_ACCESS_LINK>https://example.com/access</REMS_ACCESS_LINK>
			</LANDING_PAGE>
		</LANDING_PAGE_SET>
	`)

	frontMatter, err := markdownWriter(xmlData, "sample_dataset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []string{
		"header: Sample cohort header",
		"dataset_title: Example synthetic dataset",
		"number_of_cases: \"42\"",
		"rems_access_link: https://example.com/access",
		"sample_images:",
		"keywords:",
		"sample",
		"synthetic",
		"Research use only",
		"Example policy text:",
		"This is a multiline policy summary.",
	}

	for _, want := range checks {
		if !strings.Contains(frontMatter, want) {
			t.Fatalf("expected %q in front matter, got: %s", want, frontMatter)
		}
	}
}
