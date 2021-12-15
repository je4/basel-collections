package files

import "embed"

//go:embed static/js/base.js
//go:embed static/css/basel-collections.css
//go:embed static/css/basel-collections.css.map
var StaticFS embed.FS

//go:embed template/root.gohtml
var RootTemplate string

//go:embed template/detail.gohtml
var DetailTemplate string
