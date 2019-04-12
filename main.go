package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	config "./config"
	handler "./handler"
)

func main() {
	dbName := os.Getenv("ENV.DB_NAME")
	dbPass := os.Getenv("ENV.DB_PASS")
	dbHost := os.Getenv("ENV.DB_HOST")
	dbPort := 5432

	connection, error := config.ConnectSQL(dbHost, dbPort, dbName, dbPass, dbName)
	if error != nil {
		fmt.Println(error)
		os.Exit(-1)
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	postHandler := handler.NewPostHandler(connection)
	router.Route("/", func(r chi.Router) {
		r.Mount("/posts", postRouter(postHandler))
	})

	fmt.Println("Server listen at :8080")
	http.ListenAndServe(":8080", router)
}

func postRouter(postHandler *handler.PostHandler) http.Handler {
	router := chi.NewRouter()
	router.Get("/", postHandler.Fetch)
	router.Get("/{id:[0-9]+}", postHandler.GetByID)
	router.Post("/", postHandler.Create)
	router.Put("/{id:[0-9]+}", postHandler.Update)
	router.Delete("/{id:[0-9]+}", postHandler.Delete)

	return router
}
