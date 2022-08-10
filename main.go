package main

import (
	"fmt"
	"log"
	"mygram/database"
	"mygram/handler"
	"mygram/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	handler.ParseConfig()
	sql := database.NewSqlConnection(handler.GetConnectionString())
	database.SqlDatabase = sql
	defer sql.CloseConnection()

	r := mux.NewRouter()
	handler.InstallUsersHandler(r)
	handler.InstallPhotosHandler(r)
	handler.InstallCommentHandler(r)
	handler.InstallSocialMediaHandler(r)
	r.Use(middleware.SecureMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf(" http://%s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
