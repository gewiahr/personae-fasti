package dto

type PlayerBrief struct {
	ExtID    string `json:"ext"`
	Username string `json:"username"`
}

type PersonalNote struct {
	PersonalNote string `json:"personalNote"`
}
