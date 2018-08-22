package sett

type InvUser struct {
	Id       string
	Name     string
	Number   string
	Geo      [2]string
	Password string
	NeedHelp bool
}

type VolUser struct {
	Name        string
	Number      string
	Geo         [2]string // ["latitude","longitude"]
	Password    string
	CanHelp     bool
	GoodReviews int
	BadReviews  int
}
