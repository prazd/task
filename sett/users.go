package sett

type InvUser struct {
	Id       string
	Name     string
	Number   string
	Password string
}

type VolUser struct {
	Name     string
	Number   string
	Geo      [2]string // ["latitude","longitude"]
	Password string
}
