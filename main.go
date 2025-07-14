package main

import (
	"log"
	"net/http"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200 ok"))
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("что-то не так")
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic caught: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/panic", recoverMiddleware(http.HandlerFunc(panicHandler)))
	http.Handle("/ok", recoverMiddleware(http.HandlerFunc(okHandler)))
	http.ListenAndServe(":8080", nil)
}
