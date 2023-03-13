package service

const darkorchid = "#7b458d"
const darkorchid1 = "#69008B"
const darkgoldenrod = "#F29F05"
const deeppink = "#EC339B"
const turquoise = "#278B9B"
const black = "black"
const dark = "#1a1a1a"
const white = "white"

const fontFamilySansSerif = "'IBM Plex Mono', monospace"
const fontFamilySerif = "'IBM Plex Serif', serif"
const fontFamilySansSerifCondensed = "'IBM Plex Sans Condensed', sans-serif"

const BoxCell = "cell"
const BoxBlank = "blank"
const BoxNews = "news"
const BoxImpressum = "impressum-page"

const top = "top"
const bottom = "bottom"

var SCHEMES = map[int]map[string]string{
	1: {"color": deeppink, "background-color": darkorchid1, "font-family": fontFamilySansSerifCondensed, "vertical-align": "bottom"},
	2: {"color": darkgoldenrod, "background-color": deeppink, "font-family": fontFamilySansSerif, "font-style": "italic", "vertical-align": "bottom"},
	3: {"color": white, "background-color": dark, "font-family": fontFamilySerif, "vertical-align": "top"},
	4: {"color": black, "background-color": turquoise, "font-family": fontFamilySerif, "vertical-align": "top"},
	5: {"color": darkgoldenrod, "background-color": darkorchid, "font-family": fontFamilySansSerif, "font-style": "italic", "vertical-align": "top"},
	6: {"color": darkorchid, "background-color": darkgoldenrod, "font-family": fontFamilySansSerif, "font-style": "italic", "vertical-align": "bottom"},
}

var IMPRESSUM = map[string]string{"color": black, "background-color": darkgoldenrod, "font-family": fontFamilySansSerif, "font-style": "italic", "vertical-align": "middle"}

var NEWSSMALL = []Grid{
	{Id: 0, Left: 1, Cols: 4, Top: 2, Rows: 3, Type: BoxNews, Scheme: SCHEMES[3], VAlign: top},   // 1
	{Id: 0, Left: 5, Cols: 4, Top: 3, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: top},   // 1
	{Id: 0, Left: 1, Cols: 1, Top: 5, Rows: 3, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top},  // 1
	{Id: 0, Left: 2, Cols: 3, Top: 5, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: top},   // 1
	{Id: 0, Left: 5, Cols: 4, Top: 4, Rows: 3, Type: BoxCell, Scheme: SCHEMES[6], VAlign: top},   // 1
	{Id: 0, Left: 1, Cols: 4, Top: 6, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: top},   // 1
	{Id: 0, Left: 5, Cols: 1, Top: 7, Rows: 2, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top},  // 1
	{Id: 0, Left: 6, Cols: 3, Top: 7, Rows: 4, Type: BoxCell, Scheme: SCHEMES[1], VAlign: top},   // 1
	{Id: 0, Left: 1, Cols: 2, Top: 9, Rows: 3, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top},  // 1
	{Id: 0, Left: 3, Cols: 3, Top: 9, Rows: 3, Type: BoxCell, Scheme: SCHEMES[6], VAlign: top},   // 1
	{Id: 0, Left: 6, Cols: 3, Top: 11, Rows: 2, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top}, // 1
	{Id: 0, Left: 1, Cols: 5, Top: 12, Rows: 3, Type: BoxCell, Scheme: SCHEMES[4], VAlign: top},  // 1
	{Id: 0, Left: 3, Cols: 3, Top: 13, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: top},  // 1
	{Id: 0, Left: 5, Cols: 1, Top: 14, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top}, // 1
}

var NEWSLARGE = []Grid{
	{Id: 0, Left: 1, Cols: 4, Top: 2, Rows: 7, Type: BoxNews, Scheme: SCHEMES[3], VAlign: top},   // 1
	{Id: 0, Left: 5, Cols: 4, Top: 2, Rows: 7, Type: BoxCell, Scheme: SCHEMES[6], VAlign: top},   // 2
	{Id: 0, Left: 9, Cols: 4, Top: 4, Rows: 7, Type: BoxCell, Scheme: SCHEMES[5], VAlign: top},   // 3
	{Id: 0, Left: 1, Cols: 5, Top: 9, Rows: 6, Type: BoxCell, Scheme: SCHEMES[1], VAlign: top},   // 4
	{Id: 0, Left: 6, Cols: 3, Top: 9, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top},  // 5
	{Id: 0, Left: 9, Cols: 4, Top: 9, Rows: 6, Type: BoxCell, Scheme: SCHEMES[4], VAlign: top},   // 6
	{Id: 0, Left: 1, Cols: 2, Top: 15, Rows: 6, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: top}, // 7
	{Id: 0, Left: 3, Cols: 6, Top: 15, Rows: 6, Type: BoxCell, Scheme: SCHEMES[1], VAlign: top},  // 8
	{Id: 0, Left: 14, Cols: 4, Top: 4, Rows: 3, Type: BoxBlank, Scheme: SCHEMES[6], VAlign: top}, // 9
}

var GRIDSMALL = []Grid{
	{Id: 0, Left: 1, Cols: 3, Top: 2, Rows: 5, Type: BoxCell, Scheme: SCHEMES[1], VAlign: bottom},   // 1
	{Id: 0, Left: 4, Cols: 5, Top: 2, Rows: 3, Type: BoxNews, Scheme: SCHEMES[3], VAlign: bottom},   // 1
	{Id: 0, Left: 4, Cols: 5, Top: 5, Rows: 3, Type: BoxCell, Scheme: SCHEMES[4], VAlign: top},      // 1
	{Id: 0, Left: 1, Cols: 3, Top: 7, Rows: 5, Type: BoxCell, Scheme: SCHEMES[2], VAlign: bottom},   // 1
	{Id: 0, Left: 4, Cols: 5, Top: 8, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 1
	{Id: 0, Left: 4, Cols: 5, Top: 9, Rows: 4, Type: BoxCell, Scheme: SCHEMES[1], VAlign: bottom},   // 1
	{Id: 0, Left: 1, Cols: 3, Top: 12, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom}, // 1
	{Id: 0, Left: 1, Cols: 5, Top: 13, Rows: 4, Type: BoxCell, Scheme: SCHEMES[6], VAlign: bottom},  // 1
	{Id: 0, Left: 6, Cols: 3, Top: 13, Rows: 4, Type: BoxCell, Scheme: SCHEMES[2], VAlign: bottom},  // 1
}

var GRIDLARGE = []Grid{
	{Id: 0, Left: 1, Cols: 2, Top: 2, Rows: 3, Type: BoxCell, Scheme: SCHEMES[1], VAlign: bottom},   // 1
	{Id: 0, Left: 3, Cols: 2, Top: 2, Rows: 3, Type: BoxNews, Scheme: SCHEMES[3], VAlign: bottom},   // 2
	{Id: 0, Left: 5, Cols: 3, Top: 2, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 3
	{Id: 0, Left: 8, Cols: 2, Top: 2, Rows: 5, Type: BoxCell, Scheme: SCHEMES[6], VAlign: bottom},   // 4
	{Id: 0, Left: 10, Cols: 3, Top: -1, Rows: 4, Type: BoxCell, Scheme: SCHEMES[4], VAlign: bottom}, // 5
	{Id: 0, Left: 5, Cols: 3, Top: 3, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: bottom},   // 6
	{Id: 0, Left: 10, Cols: 3, Top: 3, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom}, // 7
	{Id: 0, Left: 10, Cols: 3, Top: 4, Rows: 3, Type: BoxCell, Scheme: SCHEMES[5], VAlign: bottom},  // 8
	{Id: 0, Left: 1, Cols: 2, Top: 5, Rows: 4, Type: BoxCell, Scheme: SCHEMES[6], VAlign: bottom},   // 9
	{Id: 0, Left: 3, Cols: 2, Top: 5, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 10
	{Id: 0, Left: 3, Cols: 4, Top: 6, Rows: 3, Type: BoxCell, Scheme: SCHEMES[4], VAlign: bottom},   // 11
	{Id: 0, Left: 7, Cols: 1, Top: 6, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 12
	{Id: 0, Left: 7, Cols: 4, Top: 7, Rows: 3, Type: BoxCell, Scheme: SCHEMES[1], VAlign: bottom},   // 13
	{Id: 0, Left: 11, Cols: 2, Top: 7, Rows: 6, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom}, // 14
	{Id: 0, Left: 1, Cols: 1, Top: 9, Rows: 4, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 15
	{Id: 0, Left: 2, Cols: 2, Top: 9, Rows: 4, Type: BoxCell, Scheme: SCHEMES[5], VAlign: bottom},   // 16
	{Id: 0, Left: 4, Cols: 2, Top: 9, Rows: 4, Type: BoxCell, Scheme: SCHEMES[6], VAlign: bottom},   // 17
	{Id: 0, Left: 6, Cols: 1, Top: 9, Rows: 1, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom},  // 18
	{Id: 0, Left: 6, Cols: 3, Top: 10, Rows: 3, Type: BoxCell, Scheme: SCHEMES[2], VAlign: bottom},  // 19
	{Id: 0, Left: 9, Cols: 1, Top: 10, Rows: 3, Type: BoxBlank, Scheme: SCHEMES[3], VAlign: bottom}, // 20
}

func buildGrid(template []Grid, collections []Content) ([]Grid, int64) {
	var grid []Grid

	var entries int64 = int64(len(collections))
	var boxes int64 = int64(len(template))
	var count int64 = 0
	var nr int64 = 0
	var firstRow int64 = 2
	var lastRow int64 = 2
	var height int64 = 0

	for _, t := range template {
		end := t.Top + t.Rows
		if end > height {
			height = end
		}
	}

	height -= firstRow
	for count < entries && nr < 300 {
		var h Grid = template[nr%boxes]
		var content Content = &StaticContent{}
		var id int64 = 0
		var iter = (nr - (nr % boxes)) / boxes
		var top = h.Top + iter*height
		h.Top = top

		if top < 0 {
			h.Top = 2
			h.Rows = top + firstRow
		} else {
			switch template[nr%boxes].Type {
			case BoxBlank:
			case BoxNews:
				content = &StaticContent{title: "News"}
			default:
				content = collections[count]
				count++
			}
		}
		h.Content = content
		h.Id = id
		n := h.Top + h.Rows
		if lastRow < n {
			lastRow = n
		}
		grid = append(grid, h)
		nr++
	}

	for i := int64(0); i < 20; i++ {
		var h = template[(nr+i)%boxes]
		var iter = ((nr + i) - ((nr + i) % boxes)) / boxes
		var top = h.Top + iter*height
		if top > lastRow {
			continue
		}
		if h.Rows+top > lastRow {
			h.Rows = lastRow - top
			if h.Rows <= 0 {
				continue
			}
		}
		h.Top = top
		h.Content = &StaticContent{}
		if h.Type == BoxNews {
			h.Content = &StaticContent{title: "News"}
		}
		n := h.Top + h.Rows
		if lastRow < n {
			lastRow = n
		}
		grid = append(grid, h)
	}
	return grid, lastRow
}
