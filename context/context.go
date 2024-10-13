package context

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// Context - контекст приложения
type Context struct {
	Port           string        `json:"port"`
	SecretKey      string        `json:"secret_key"`
	ExpireDuration time.Duration `json:"expire_duration"`
	DBHost         string        `json:"db_host"`
	DBPort         string        `json:"db_port"`
	DBUser         string        `json:"db_user"`
	DBPassword     string        `json:"db_password"`
	DBName         string        `json:"db_name"`
	EmailUser      string        `json:"email_user"`
	EmailPassword  string        `json:"email_password"`
}

var appContext *Context

// InitContext - инициализирует контекст приложения
func InitContext() (Context, error) {
	if appContext == nil {
		appContext = &Context{}
		if err := godotenv.Load(); err != nil {
			return Context{}, errors.New(".env файл не найден")
		}
		port, exists := os.LookupEnv("PORT")
		if !exists {
			return Context{}, errors.New("не найден порт в конфигурационном файле")
		}
		appContext.Port = port

		secretKey, exists := os.LookupEnv("SECRET_KEY")
		if !exists {
			return Context{}, errors.New("не найден ключ для авторизации в конфигурационном файле")
		}
		appContext.SecretKey = secretKey

		expireDurationStr, exists := os.LookupEnv("EXPIRE_DURATION")
		if !exists {
			return Context{}, errors.New("не найдено время существованию токена в конфигурационном файле")
		}
		expireDurationInt, err := strconv.Atoi(expireDurationStr)
		if err != nil {
			return Context{}, fmt.Errorf("некорректный формат времени существования токена: %v", err)
		}
		appContext.ExpireDuration = time.Duration(expireDurationInt)

		dbHost, exists := os.LookupEnv("DB_HOST")
		if !exists {
			return Context{}, errors.New("не найден хост базы данных в конфигурационном файле")
		}
		appContext.DBHost = dbHost

		dbPort, exists := os.LookupEnv("DB_PORT")
		if !exists {
			return Context{}, errors.New("не найден порт базы данных в конфигурационном файле")
		}
		appContext.DBPort = dbPort

		dbUser, exists := os.LookupEnv("DB_USER")
		if !exists {
			return Context{}, errors.New("не найдено имя пользователя базы данных в конфигурационном файле")
		}
		appContext.DBUser = dbUser

		dbPassword, exists := os.LookupEnv("DB_PASSWORD")
		if !exists {
			return Context{}, errors.New("не найдено пароль от базы данных в конфигурационном файле")
		}
		appContext.DBPassword = dbPassword

		dbName, exists := os.LookupEnv("DB_NAME")
		if !exists {
			return Context{}, errors.New("не найдено имя базы данных в конфигурационном файле")
		}
		appContext.DBName = dbName

		emailUser, exists := os.LookupEnv("EMAIL_USER")
		if !exists {
			return Context{}, errors.New("не найдена почта пользователя в конфигурационном файле")
		}
		appContext.EmailUser = emailUser

		emailPassword, exists := os.LookupEnv("EMAIL_PASSWORD")
		if !exists {
			return Context{}, errors.New("не найден пароль от почты в конфигурационном файле")
		}
		appContext.EmailPassword = emailPassword
	}
	return *appContext, nil
}

// GetApiContext - геттер контекст приложения
func GetApiContext() Context {
	return *appContext
}
