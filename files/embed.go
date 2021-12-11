package files

import "embed"

//go:embed static/*
var StaticFS embed.FS

//go:embed template/*
var TemplateFS embed.FS

var TemplateFiles = map[string]string{}
