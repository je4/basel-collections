package files

import "embed"

//go:embed static/js/base.js
//go:embed static/css/basel-collections.css
//go:embed static/css/basel-collections.css.map
var StaticFS embed.FS

//go:embed template/*
var TemplateFS embed.FS

var TemplateFiles = map[string]string{
	"root": "template/root.gohtml",
}
