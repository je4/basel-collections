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
var TemplateFS embed.FS

//go:embed template/header_content.inc.gohtml
var HeaderContentIncTemplate string

//go:embed template/collections.gohtml
var CollectionsTemplate string

//go:embed template/collection.gohtml
var CollectionTemplate string

//go:embed template/news.gohtml
var NewsTemplate string
