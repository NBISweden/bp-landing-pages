---
hidden: no
datacite:
  "@context": "http://schema.org"
  "@type": "Dataset"
  "@id": "https://doi.org/10.123/70x8wi7vq"
  name: "Mock data 38"
  about: "Pathology"
  url: "https:/.Mock"
  author:
  - name: "Person 1"
    #"@id": # FIXME: missing info
    "@type": "Person"
  - name: "Person 2"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  - name: "Person 3"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  - name: "Person 4"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  - name: "Person 5"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  publisher:
    "@type": "Organization"
    name: "BP"
  copyrightYear: 2020
  copyrightHolder:
  - name: "University 2"
    url: "https://uni2.se/"
    "@type": "Organization"
  - name: "person 1"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  provider:
  - name: "Person 2"
    email: "person1@uni1.se"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  - name: "person 3"
    email: "person3@uni1.se"
    #"@id": "" # FIXME: missing info
    "@type": "Person"        
  - name: "Person 4"
    email: "Person4@uni1.se"
    "@id": "https://orcid.org/0"
    "@type": "Person"
  - name: "BP"
    email: "BP@mock.com"
    "@id": "https://datahub.bp.mock.se"
    "@type": "Organization"
  dateCreated: "2023-02-05"
  datePublished: "2023-02-07"
  dateModified: "2023-02-09"
  keywords: "Prostate Cancer, Pathology"
  version: "1.0.0"
  description: |
    2000 scanned prostate biopsies from 120 patient cases.
  license:
  - name: "Controlled access"
    id: "https://datahub.bp.mock.se/mock#controlled-access"
    "@type": "CreativeWork"
    abstract: |
      Free for use in legal and ethical medical diagnostics research.
  - name: "Big picture license"
    id: "https://"
  citation:
other:
  shortName: "MD38"
  status: "Completed"
  annotation: |

  countries-shared: # Fixme: insert link to paper in annotation text above when published.
  - "SE"
  organ:
  - name: "Prostate"
    sctid: 41216001 # SNOMED-CT https://termbrowser.nhs.uk/?perspective=full&conceptId1=%s
  age-span: "39 - 79 years, median 67 years"
  bytes: 492600000000 # ~160 GB
  numberOfScans: 2611
  caseLevelAnnotaions: "Treatment decision, Gleason, PIRADS, PSA, Prostate volume, clinical T-stage"
  imageLevelAnnotations: "Connection to use"
  numberOfAnnotations:  # 652 (x8 echo times) + 40 (x8 echo times)
  resolution: "0.2204 microns per pixel (40x)"
  modality:
  - "SM"
  scanner:
  - "Hamamatsu NanoZoomer S60"
  stain: "H&E"
  phase:
  image: "img/m_p_p/lymph-nodes-to-scale.jpeg"
  exampleImage:
  - title: "Example biopsy downsampled 8x"
    url: "/assets/images/mock_datasets/bp/m_p_p/lymph-nodes-to-scale.jpeg"
    thumbnail-url: "/assets/images/mock_datasets/bp/m_p_p/lymph-nodes-to-scale-thumbnail.jpeg"
  - title: "Example biopsy downsampled 8x"
    url: "/assets/images/mock_datasets/bp/m_p_p/b.jpg"
    thumbnail-url: "/assets/images/mock_datasets/bp/m_p_p/b_thumb.jpg"
  - title: "Example patch in full resolution"
    url: "/assets/images/mock_datasets/bp/m_p_p/pa_fullres.jpg"
    thumbnail-url: "/assets/images/mock_datasets/bp/m_p_p/pa_fullres_thumb.jpg"
  - title: "Macro image with annotation"
    url: "/assets/images/mock_datasets/bp/m_p_p/macro_a.jpg"
    thumbnail-url: "/assets/images/mock_datasets/bp/m_p_p/macro_a_thumb.jpg"
access: Indirect
---


## File formats

Images: Pyramidal .tif

Metadata: .csv

## License






Copyright
{{ page.datacite.copyrightYear }}
{{ page.datacite.copyrightHolder | map: "name" |  join: ", " }}

Donec nibh elit, tincidunt eget vulputate ac, consequat in ipsum. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Cras interdum enim at libero sodales finibus. Donec vitae lorem quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Phasellus eget orci vulputate quam porttitor dignissim. Maecenas commodo, nibh et iaculis blandit, leo arcu scelerisque tortor, id porttitor metus tortor commodo quam. Vivamus a dapibus ipsum, eu euismod elit.

{{ page.datacite.author | map: "name" | array_to_sentence_string }}
({{ page.datacite.datePublished | date: "%Y" }})
{{ page.datacite.name }}
[doi:{{ page.datacite['@id'] | remove: "https://doi.org/" }}]({{ page.datacite["@id"] }}).

 Integer tempus ullamcorper quam, eu consequat ante pulvinar ac. Nulla in ipsum mattis, congue orci vitae, rhoncus lacus. Aliquam nunc turpis, tincidunt in nibh eu, bibendum luctus diam. Praesent varius orci erat, nec sollicitudin tellus venenatis vel. Morbi quis ullamcorper arcu, quis ultricies odio.
