package home

import (
	"math"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/tools/datasize"
)

type HomePageTmpl struct {
	User      *users.User
	MemoryBar ValueBar
	SwapBar   ValueBar
}

func (t *HomePageTmpl) Template() string { return "home/page_home" }

type ValueBar struct {
	Label string
	Value datasize.ByteSize
	Total datasize.ByteSize
}

func (v *ValueBar) Percentage() int {
	return int(math.Round(float64(v.Value) / float64(v.Total) * 100))
}
