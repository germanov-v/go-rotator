package main

import (
	"fmt"
	"github.com/germanov-v/go-rotator/internal/config"
	"github.com/germanov-v/go-rotator/internal/repository/postgres"
	"github.com/germanov-v/go-rotator/internal/service"
	"github.com/germanov-v/go-rotator/internal/transport/http_handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("rotator")
	config, _ := config.ReadConsoleConfigParameter()

	repo, err := postgres.NewPostgresRepo(config.DbBaseConfig.ConnectionString)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	r := mux.NewRouter()
	service := service.NewRotationService(repo)

	r.HandleFunc("/slots/{slot}/banners", http_handler.AddBannerHandler(repo)).Methods("POST")
	r.HandleFunc("/slots/{slot}/banners/{banner}", http_handler.RemoveBannerHandler(repo)).Methods("DELETE")
	r.HandleFunc("/slots/{slot}/rotate", http_handler.RotateBannerHandler(service)).Methods("GET")
	r.HandleFunc("/slots/{slot}/stats/click", http_handler.RecordClickHandler(repo)).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server listening on :! " + strconv.Itoa(config.ServerConfig.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.ServerConfig.Port), nil))
}
