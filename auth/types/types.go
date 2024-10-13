package types

import (
	"github.com/dgrijalva/jwt-go/v4"
)

// Claims - расширенные claims для jwt токена
type Claims struct {
	jwt.StandardClaims
	GUID string `json:"guid"`
	Ip   string `json:"ip"`
}

// UserSessionIdentity - информация о пользователе
type UserSessionIdentity struct {
	GUID         string // guid пользователя
	Ip           string // ip адрес пользователя
	RefreshToken string // refresh_token пользователя
	Email        string // email пользователя
}
