package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"pmj-server/controller"
)

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://localhost:8081")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Redmine-API-Key, X-CSRF-Token")

		fs.ServeHTTP(w, r)
	}
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	fs := http.FileServer(http.Dir("view/build/"))

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		fmt.Printf("読み込み出来ませんでした: %v", errEnv)
	}
	http.Handle("/", cors(fs))
	http.HandleFunc("/user", controller.UserController{}.Execute)
	http.HandleFunc("/config", controller.ConfigController{}.Execute)

	err := server.ListenAndServeTLS("data_cert.pem", "data_key.pem")
	if err != nil {
		return
	}
}
