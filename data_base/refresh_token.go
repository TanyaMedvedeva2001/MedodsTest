package data_base

import (
	"authTestMedods/auth/types"
	"fmt"
)

// WriteRefreshToken - записывает данных о сессии (включая токен) в бд
func (db *DBManager) WriteRefreshToken(userRefrInfo types.UserSessionIdentity) error {
	query := `INSERT INTO session_info (guid, ip, refresh_token, date_added, date_update) VALUES ($1, $2, $3, CURRENT_DATE, CURRENT_DATE)`
	_, err := db.DB.Exec(query, userRefrInfo.GUID, userRefrInfo.Ip, userRefrInfo.RefreshToken)
	return err
}

// GetGUIDFromToken - получает guid пользователя по refresh токену
func (db *DBManager) GetGUIDFromToken(token string) (string, error) {
	row := db.DB.QueryRow(`SELECT guid FROM session_info WHERE refresh_token = $1`, token)
	var guid string
	err := row.Scan(&guid)
	if err != nil {
		return "", fmt.Errorf("не удалось получить guid по refresh_token: %v", err)
	}
	return guid, nil
}

// UpdateRefreshToken - обновляет refresh токен в бд
func (db *DBManager) UpdateRefreshToken(userInfo types.UserSessionIdentity) (int64, error) {
	result, err := db.DB.Exec(`UPDATE session_info set refresh_token = $1, date_update = CURRENT_DATE where guid = $2`,
		userInfo.RefreshToken, userInfo.GUID)
	if err != nil {
		return 0, fmt.Errorf("не удалось обновить refresh_token в базе данных: %v", err)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество измененных строк: %v", err)
	}
	return rowAffected, nil
}

// GetIp - возврашает ip пользователя по его guid
func (db *DBManager) GetIp(guid string) (string, error) {
	row := db.DB.QueryRow(`SELECT ip FROM session_info WHERE guid = $1`, guid)
	var ip string
	err := row.Scan(&ip)
	if err != nil {
		return "", fmt.Errorf("не удалось получить ip сессии: %v", err)
	}
	return ip, nil
}

// UpdateIp - обновляет Ip пользователя
func (db *DBManager) UpdateIp(guid string, ip string) (int64, error) {
	result, err := db.DB.Exec("UPDATE session_info set ip = $1 where guid = $2",
		ip, guid)
	if err != nil {
		return 0, fmt.Errorf("не удалось обновить ip в базе данных: %v", err)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество измененных строк: %v", err)
	}
	return rowAffected, nil
}
