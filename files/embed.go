package files

import "embed"

//go:embed static/js/base.js
//go:embed static/css/basel-collections.css
//go:embed static/css/basel-collections.css.map
var StaticFS embed.FS

//go:embed template/collections.gohtml
var CollectionsTemplate string

//go:embed template/collection.gohtml
var CollectionTemplate string

//go:embed template/news.gohtml
var NewsTemplate string
