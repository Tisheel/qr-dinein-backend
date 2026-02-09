package model

type LoginRequest struct {
	Username string `json:"username"`
	Pin      string `json:"pin"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	Staff     *Staff `json:"staff"`
}

type SuperuserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SuperuserLoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	Role      string `json:"role"`
	Username  string `json:"username"`
}
