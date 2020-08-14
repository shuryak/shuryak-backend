package main

import (
	"flag"
	"fmt"
	"github.com/shuryak/shuryak-backend/internal/handlers/articles"
	"github.com/shuryak/shuryak-backend/internal/handlers/users"
	"github.com/shuryak/shuryak-backend/internal/middleware"
	"github.com/shuryak/shuryak-backend/internal/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func handleRequests() http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.HeadersMiddleware)
	router.HandleFunc("/api/articles.create", middleware.IsAuthMiddleware(articles.CreateHandler))
	router.HandleFunc("/api/articles.findOne", articles.FindOneHandler)
	router.HandleFunc("/api/articles.findMany", articles.FindManyHandler)
	router.HandleFunc("/api/articles.getById", articles.GetByCustomIdHandler)
	router.HandleFunc("/api/articles.getList", articles.GetListHandler)
	router.HandleFunc("/api/users.register", users.CreateHandler)
	router.HandleFunc("/api/users.login", users.LoginHandler)
	router.HandleFunc("/api/users.refreshTokenPair", users.RefreshTokenPairHandler)

	http.Handle("/", router)

	return cors.Default().Handler(router)
}

func main() {
	profile := flag.String("profile", "debug", "Configuration profile selection")
	flag.Parse()

	var config *utils.ProfileType

	if *profile == "debug" {
		config = utils.Configuration.Debug
	} else if *profile == "release" {
		config = utils.Configuration.Release
	} else {
		log.Fatal("Bad profile!")
	}

	utils.OpenMongo("mongodb://localhost:27017")
	defer utils.CloseMongo()

	fmt.Println("Server is running on", *config.ServerPort, "port!")
	err := http.ListenAndServe(":"+*config.ServerPort, handleRequests())
	if err != nil {
		log.Fatal("Internal error!")
	}
}
