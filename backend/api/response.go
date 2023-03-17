package api

// TODO
type CitationSupplemental struct{}

type Link struct {
	DisplayName  string `json:"display_name"`
	Url          string `json:"url"`
	CoverageText string `json:"coverage_text"`
}

type Record struct {
	CitationSupplemental CitationSupplemental `json:"citation_supplemental"`
	Links                []Link               `json:"links"`
}

type Response struct {
	Errors  []string `json:"errors"`
	Found   bool     `json:"found"`
	Records []Record `json:"records"`
}
