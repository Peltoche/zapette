package sysstats

type SysstatsPageTmpl struct {
	GraphData *Graph
}

func (t *SysstatsPageTmpl) Template() string { return "sysstats/page_list" }

type Dataset struct {
	Label       string     `json:"label"`
	Data        []*float64 `json:"data"`
	ShowLine    bool       `json:"showLine"`
	BorderColor string     `json:"borderColor"`
	SteppedLine bool       `json:"steppedLine"`
	BorderWidth int        `json:"borderWidth"`
	PointRadius int        `json:"pointRadius"`
}

type Data struct {
	Labels   []*string `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Graph struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}
