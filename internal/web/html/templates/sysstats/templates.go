package sysstats

import (
	"math"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/tools/datasize"
)

type SysstatsPageTmpl struct {
	User      *users.User
	MemoryBar ValueBar
	SwapBar   ValueBar
}

func (t *SysstatsPageTmpl) Template() string { return "sysstats/page_list" }

type ValueBar struct {
	Label string
	Value datasize.ByteSize
	Total datasize.ByteSize
}

func (v *ValueBar) Percentage() int {
	return int(math.Round(float64(v.Value) / float64(v.Total) * 100))
}
