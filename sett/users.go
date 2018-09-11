package sett

type InvUser struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Geo      [2]string `json:"geo"`
	Password string    `json:"password"`
	State    int       `json:"state"`
	ConID    int       `json:"conid"`
	Online   bool      `json:"online"`
}

type VolUser struct {
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Geo         [2]string `json:"geo"`
	Password    string    `json:"password"`
	GoodReviews int       `json:"goodreviews"`
	BadReviews  int       `json:"badreviews"`
	State       int       `json:"state"`
	Online      bool      `json:"online"`
}
