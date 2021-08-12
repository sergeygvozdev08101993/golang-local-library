package main

import (
	"context"
	"net/http"

	"github.com/sergeygvozdev08101993/golang-local-library/internal/app"
	"github.com/sergeygvozdev08101993/golang-local-library/pkg/service/handlers"
)

func main() {

	var err error

	app.InitLogger()

	ctx := context.Background()
	app.ClientDB, err = app.ConnectAndPingToClientDB()
	if err != nil {
		app.ErrLog.Fatalf("failed to connect with client database: %v", err)
	}
	defer app.ClientDB.Disconnect(ctx)

	router := http.NewServeMux()
	router.HandleFunc("/", handlers.Index)
	router.HandleFunc("/catalog", handlers.Catalog)
	router.HandleFunc("/catalog/books", handlers.ListBooks)
	router.HandleFunc("/catalog/authors", handlers.ListAuthors)
	router.HandleFunc("/catalog/genres", handlers.ListGenres)
	router.HandleFunc("/catalog/bookinstances", handlers.ListBookInstances)

	router.HandleFunc("/catalog/author/create", handlers.CreateAuthor)
	router.HandleFunc("/catalog/book/create", handlers.CreateBook)
	router.HandleFunc("/catalog/genre/create", handlers.CreateGenre)
	router.HandleFunc("/catalog/bookinstance/create", handlers.CreateBookInstance)

	app.InfoLog.Println("server is running...")
	if err = http.ListenAndServe(":7000", router); err != nil {
		app.ErrLog.Fatalf("server is stopped: %v", err)
	}
}
