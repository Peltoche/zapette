package sysstats

import (
	"github.com/Peltoche/zapette/internal/service/sysstats"
	"github.com/Peltoche/zapette/internal/service/users"
)

type SysstatsPageTmpl struct {
	Latest      sysstats.Stats
	User        *users.User
	Labels      []*string
	MemoryUsed  []*float64
	MemoryTotal []*float64
	SwapUsed    []*float64
	CacheBuffer []*float64
}

func (t *SysstatsPageTmpl) Template() string { return "sysstats/page_list" }
