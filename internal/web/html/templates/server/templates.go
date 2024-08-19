package server

import (
	"github.com/Peltoche/zapette/internal/service/sysstats"
)

type DetailsPageTmpl struct {
	Stats *sysstats.Stats
}

func (t *DetailsPageTmpl) Template() string { return "server/page_details" }
