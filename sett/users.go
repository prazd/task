package sett

type InvUser struct {
	Id       string
	Name     string
	Phone    string
	Geo      [2]string // ["latitude","longitude"]
	Password string
	// NeedHelp bool
	// Busy     bool
	State int
}

type VolUser struct {
	Name     string
	Phone    string
	Geo      [2]string
	Password string
	// CanHelp     bool
	GoodReviews int
	BadReviews  int
	// Busy        bool
	State int
}
