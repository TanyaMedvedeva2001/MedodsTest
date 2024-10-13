package data_base

import (
	app_context "authTestMedods/context"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DBManager - структура менеджера базы данных
type DBManager struct {
	DB *sql.DB
}

var dbManager *DBManager

// GetDBManager - геттер менеджера базы данных
func GetDBManager() *DBManager {
	return dbManager
}

// InitDBManager - инициализация менеджера базы данных
func InitDBManager(ctx context.Context, appContext app_context.Context) error {
	if dbManager == nil {
		dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			appContext.DBHost, appContext.DBPort, appContext.DBUser, appContext.DBPassword, appContext.DBName)
		db, err := sql.Open("postgres", dbInfo)
		if err != nil {
			return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
		}
		dbManager = &DBManager{DB: db}
		go func(ctx context.Context, db *sql.DB) {
			<-ctx.Done()
			db.Close()
		}(ctx, db)
	}
	return nil
}
