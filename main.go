package main

import (
	"flag"
	"fmt"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/articles"
	"github.com/shuryak/shuryak-backend/internal/middleware"
	"github.com/shuryak/shuryak-backend/internal/users"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter()

	router.Use(middleware.HeadersMiddleware)
	router.HandleFunc("/api/articles.create", middleware.IsAuthMiddleware(articles.CreateHandler))
	router.HandleFunc("/api/articles.findOne", articles.FindOneHandler)
	router.HandleFunc("/api/articles.findMany", articles.FindManyHandler)
	router.HandleFunc("/api/articles.getById", articles.GetByIdHandler)
	router.HandleFunc("/api/articles.getList", articles.GetListHandler)
	router.HandleFunc("/api/users.register", users.CreateHandler)
	router.HandleFunc("/api/users.login", users.LoginHandler)

	http.Handle("/", router)
}

func main() {
	profile := flag.String("profile", "debug", "Configuration profile selection")
	flag.Parse()

	var config *internal.ProfileType

	if *profile == "debug" {
		config = internal.Configuration.Debug
	} else if *profile == "release" {
		config = internal.Configuration.Release
	} else {
		log.Fatal("Bad profile!")
	}

	handleRequests()

	internal.OpenMongo("mongodb://localhost:27017")
	defer internal.CloseMongo()

	fmt.Println("Server is running on", *config.ServerPort, "port!")
	err := http.ListenAndServe(":"+*config.ServerPort, nil)
	if err != nil {
		log.Fatal("Internal error!")
	}
}
