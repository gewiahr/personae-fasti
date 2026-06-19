package dto

type LoginRequest struct {
	Username    string `json:"username"`
	LoginSource string `json:"loginSource"`
	LoginData   string `json:"loginData"`
}

type LoginInfo struct {
	Authorization string          `json:"authorization"`
	Player        LoginPlayerInfo `json:"player"`
}

type LoginPlayerInfo struct {
	ExtID            string                   `json:"ext"`
	Username         string                   `json:"username"`
	CurrentGameExtID string                   `json:"gameExt,omitzero"`
	Settings         *LoginPlayerInfoSettings `json:"settings"`
}

type LoginPlayerInfoSettings struct {
	CouldChangeUsername bool `json:"couldChangeUsername"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
