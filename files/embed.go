package files

import "embed"

//go:embed static/js/base.js
//go:embed static/css/basel-collections.css
//go:embed static/css/basel-collections.css.map
var StaticFS embed.FS

//go:embed template/header_content.inc.gohtml
//go:embed template/collections.gohtml
//go:embed template/collection.gohtml
//go:embed template/news.gohtml
//go:embed template/impressum.gohtml
//go:embed template/information.gohtml
//go:embed template/datenschutz.gohtml
//go:embed template/kontakt.gohtml
var TemplateFS embed.FS
