package domain

// UserRepository はユーザー情報を管理するインターフェースです。
type UserRepository interface {
	RegisterUser(name string, pass string) error
	AuthenticateUser(name string, pass string) bool
	IsAccessTokenValid(token string) bool
	SelectAccessToken(username string, pass string) string
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	StatusCode  int    `json:"statusCode"`
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}
