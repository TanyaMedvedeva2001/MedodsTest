package server

import (
	"authTestMedods/auth"
	app_context "authTestMedods/context"
	"authTestMedods/server/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

// getTokens - handler, отвечающий за REST API запрос на получение токенов
func getTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	authObj := auth.GetAuthManager()
	accessToken, refreshToken, err := authObj.ReturnTokens(params["guid"], r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	jsonResponse, err := json.Marshal(types.Response{RefreshToken: refreshToken, AccessToken: accessToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonResponse)
}

// refresh - handler, отвечающий за REST API запрос на обновление токенов
func refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	authObj := auth.GetAuthManager()
	accessToken, refreshToken, err := authObj.UpdateTokens(params["refresh_token"], r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	jsonResponse, err := json.Marshal(types.Response{RefreshToken: refreshToken, AccessToken: accessToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonResponse)
}

// AppStart - запуск сервера, объявление hadler-ов
func AppStart(ctx context.Context, appContext app_context.Context) {
	httpServer := &http.Server{Addr: appContext.Port, BaseContext: func(_ net.Listener) context.Context {
		return ctx
	}}
	r := mux.NewRouter()
	r.HandleFunc("/tokens/{guid}", getTokens).Methods("GET")
	r.HandleFunc("/refresh/{refresh_token}", refresh).Methods("GET")
	httpServer.Handler = r
	g, gCtx := errgroup.WithContext(ctx)
	fmt.Printf("Сервер запущен на порту: %s\n", appContext.Port)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("Сервер остановил свою работу: %v\n", err)
	}
}
