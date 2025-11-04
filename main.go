package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kizoukun/codingtest/controller"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/mock"
)

func main() {
	r := mux.NewRouter()

	godotenv.Load()
	mock.InitDbIncremental()
	helper.InitHelper()

	// register handlers with methods
	r.HandleFunc("/api/v1/todos", controller.GetTodoController).Methods("GET")
	r.HandleFunc("/api/v1/todos", controller.AddTodoController).Methods("POST")
	r.HandleFunc("/api/v1/todos/{id}", controller.DeleteTodoController).Methods("DELETE")
	r.HandleFunc("/api/v1/todos/{id}/complete", controller.ToggleTodoController).Methods("POST")

	r.HandleFunc("/api/v1/auth/login", controller.LoginController).Methods("POST")
	r.HandleFunc("/api/v1/auth/register", controller.RegisterController).Methods("POST")

	// JSON content-type middleware (optional)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(), // remove if you don't use cookies/credentials
	)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(r)); err != nil {
		log.Fatal(err)
	}
}
