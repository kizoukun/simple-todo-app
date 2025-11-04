package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kizoukun/codingtest/controller"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/middleware"
	"github.com/kizoukun/codingtest/mock"
)

func main() {
	r := mux.NewRouter()

	godotenv.Load()
	mock.InitDbIncremental()
	helper.InitHelper()

	todos := r.PathPrefix("/api/v1").Subrouter()
	todos.Use(middleware.JWTMiddleware)

	// boards
	todos.HandleFunc("/boards", controller.GetTodoBoardController).Methods("GET")
	todos.HandleFunc("/boards", controller.AddTodoBoardController).Methods("POST")
	todos.HandleFunc("/boards/{id}", controller.DeleteTodoBoardController).Methods("DELETE")

	//team invite
	// register handlers with methods
	todos.HandleFunc("/todos/{board_id}", controller.GetTodoController).Methods("GET")
	todos.HandleFunc("/todos/{board_id}", controller.AddTodoController).Methods("POST")
	todos.HandleFunc("/todos/{board_id}/invite", controller.InviteTodoTeamController).Methods("POST")
	todos.HandleFunc("/todos/{board_id}/{id}", controller.DeleteTodoController).Methods("DELETE")
	todos.HandleFunc("/todos/{board_id}/{id}/complete", controller.ToggleTodoController).Methods("POST")

	r.HandleFunc("/api/v1/accept", controller.AcceptTodoTeamInviteController).Methods("POST")
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
