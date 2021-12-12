package service

type Grid struct {
	Id, Left, Cols, Top, Rows int64
	Type                      string
	Scheme                    map[string]string
	VAlign                    string
	Text                      string
}
