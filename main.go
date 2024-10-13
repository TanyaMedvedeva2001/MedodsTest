package main

import (
	"authTestMedods/auth"
	app_context "authTestMedods/context"
	"authTestMedods/data_base"
	"authTestMedods/mail"
	"authTestMedods/server"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// инициализируется контекст, срабатывающий при нажатии ctrl+c и останавливающий приложение
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	// инициализация контекста приложения и необходимых менеджеров
	appContext, err := app_context.InitContext()
	if err != nil {
		panic(err)
	}
	auth.InitAuthManager(appContext.SecretKey, appContext.ExpireDuration)
	err = data_base.InitDBManager(ctx, appContext)
	if err != nil {
		panic(err)
	}
	mail.InitMailManager(appContext.EmailUser, appContext.EmailPassword)
	// запуск сервера
	server.AppStart(ctx, appContext)
}
