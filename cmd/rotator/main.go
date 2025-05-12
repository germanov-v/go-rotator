package main

import (
	"fmt"
	"github.com/germanov-v/go-rotator/internal/config"
	"github.com/germanov-v/go-rotator/internal/repository/postgres"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("rotator")
	config, _ := config.ReadConsoleConfigParameter()

	repo, err := postgres.NewPostgresRepo(config.DbBaseConfig.ConnectionString)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	r := mux.NewRouter()

	//r.HandleFunc("/slots/{slot}/banners", )

	http.Handle("/", r)
}
