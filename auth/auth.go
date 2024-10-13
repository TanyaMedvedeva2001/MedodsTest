package auth

import (
	"authTestMedods/auth/types"
	"authTestMedods/data_base"
	"authTestMedods/mail"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// AuthManager - структура менеджера авторизации
type AuthManager struct {
	secretKey      string
	expireDuration time.Time
}

var auth *AuthManager

// InitAuthManager - инициализация менеджера авторизации
func InitAuthManager(secretKey string, expireDuration time.Duration) {
	if auth == nil {
		auth = &AuthManager{secretKey: secretKey, expireDuration: time.Now().Add(expireDuration * time.Minute)}
	}
	return
}

// GetAuthManager - геттера менеджера авторизации
func GetAuthManager() AuthManager {
	return *auth
}

// ReturnTokens - возвращает access и refresh токен
func (a AuthManager) ReturnTokens(guid string, ip string) (string, string, error) {
	accessToken, err := a.generateAccessToken(guid, ip)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать access token: %v", err)
	}
	refreshToken, err := a.generateRefreshToken(guid)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать refresh token: %v", err)
	}
	dbManager := data_base.GetDBManager()
	userSessionInfo := types.UserSessionIdentity{GUID: guid, Ip: ip, RefreshToken: string(refreshToken)}
	err = dbManager.WriteRefreshToken(userSessionInfo)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сохранить информацию о сессии: %v", err)
	}
	encodeRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return accessToken, encodeRefreshToken, nil
}

// getAccessToken - возращает access токен
func (a AuthManager) generateAccessToken(guid string, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &types.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(a.expireDuration),
			IssuedAt:  jwt.At(time.Now()),
		},
		GUID: guid,
		Ip:   ip,
	})
	return token.SignedString([]byte(a.secretKey))
}

// generateRefreshToken - генерирует refresh токен
func (a AuthManager) generateRefreshToken(guid string) (string, error) {
	password := []byte(guid + a.secretKey)
	hashedToken, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать refresh token: %v", err)
	}
	return string(hashedToken), err
}

// UpdateTokens - возвращает новые access и refresh токены для существующего пользователя
func (a AuthManager) UpdateTokens(refreshToken string, ip string) (string, string, error) {
	tokenDecodeByte, err := base64.StdEncoding.DecodeString(refreshToken)
	tokenDecode := string(tokenDecodeByte)
	dbManager := data_base.GetDBManager()
	guid, err := dbManager.GetGUIDFromToken(tokenDecode)
	if err != nil {
		return "", "", fmt.Errorf("не удалось получить информацию о пользователе по полученному токену: %v", err)
	}

	accessToken, err := a.generateAccessToken(guid, ip)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать access token: %v", err)
	}

	refreshToken, err = a.generateRefreshToken(guid)
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать refresh token: %v", err)
	}

	oldIp, err := dbManager.GetIp(guid)
	if err != nil {
		return "", "", fmt.Errorf("не удалось получить старый ip адрес сессии: %v", err)
	}

	userSessionInfo := types.UserSessionIdentity{GUID: guid, Ip: ip, RefreshToken: tokenDecode, Email: "tatyana-medvedeva-01@mail.ru"}
	rowAffected, err := dbManager.UpdateRefreshToken(userSessionInfo)
	if err != nil {
		return "", "", fmt.Errorf("не удалось обновить информацию о сессии: %v", err)
	}
	if rowAffected == 0 {
		return "", "", fmt.Errorf("не удалось обновить запись в базе данных")
	}

	if oldIp != ip {
		rowAffected, err = dbManager.UpdateIp(guid, ip)
		if err != nil {
			return accessToken, refreshToken, fmt.Errorf("не удалось обновить значение ip адреса: %v", err)
		}
		if rowAffected == 0 {
			return accessToken, refreshToken, fmt.Errorf("не обновились данные при новом ip адресе: %v", err)
		}
		mailManager := mail.GetMailManager()
		err = mailManager.SendMailNotification(userSessionInfo)
		if err != nil {
			return accessToken, refreshToken, fmt.Errorf("не удалось отправить письмо-оповещение о новом ip: %v", err)
		}
	}
	encodeRefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return accessToken, encodeRefreshToken, nil
}
