package home

import "github.com/Peltoche/zapette/internal/service/users"

type HomePageTmpl struct {
	User *users.User
}

func (t *HomePageTmpl) Template() string { return "home/page_home" }
