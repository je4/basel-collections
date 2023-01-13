package service

type Content interface {
	GetId() int64
	GetDate() string
	GetInst() string
	GetTitle() string
	GetUrl() string
	GetTagIDs() []int64
	GetInstitutionID() int64
}

type StaticContent struct {
	id          int64
	date        string
	institution string
	title       string
	url         string
}

func (sc *StaticContent) GetTagIDs() []int64 {
	return []int64{}
}

func (sc *StaticContent) GetInstitutionID() int64 {
	return 0
}

func (sc *StaticContent) GetId() int64     { return sc.id }
func (sc *StaticContent) GetDate() string  { return sc.date }
func (sc *StaticContent) GetInst() string  { return sc.institution }
func (sc *StaticContent) GetTitle() string { return sc.title }
func (sc *StaticContent) GetUrl() string   { return sc.url }

type Grid struct {
	Id, Left, Cols, Top, Rows int64
	Type                      string
	Scheme                    map[string]string
	VAlign                    string
	Content                   Content
}
