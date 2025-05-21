package main

import (
	"github.com/raminpz/gocourse_web/pkg/bootstrap"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/raminpz/gocourse_web/internal/user"
)

func main() {
	router := mux.NewRouter()

	_ = godotenv.Load()
	l := bootstrap.InitLoger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal("Failed to connect to database: ", err)
	}

	userRepo := user.NewRepo(l, db)

	userSrv := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
