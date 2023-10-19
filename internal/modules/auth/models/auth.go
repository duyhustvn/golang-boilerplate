package authmodel

// User kafka message
// {"username": "duyle", "password": "123"}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
