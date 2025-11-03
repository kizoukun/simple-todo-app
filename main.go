package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/usecase"
)

func main() {
	r := mux.NewRouter()

	// register handlers with methods
	r.HandleFunc("/api/v1/todos", usecase.TodoHandler).Methods("GET")
	r.HandleFunc("/api/v1/todos", usecase.AddTodoHandler).Methods("POST")
	r.HandleFunc("/api/v1/todos/{id}", usecase.DeleteTodoHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/todos/{id}/complete", usecase.ToggleTodoHandler).Methods("POST")

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
