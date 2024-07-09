package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tefsi/internal/inits"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	postgresUser := "postgres"
	postgresPassword := "password"
	postgresDB := "postgres"
	postgresHost := "localhost"
	postgresPort := "5432"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")
	defer db.Close()

	allTables, err := inits.GetAllTables(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("got all tables")

	repos, err := inits.InitRepositories(db, allTables)
	if err != nil {
		log.Fatal(err)
	}

	services := inits.InitServices(repos)
	handlers := inits.InitHandlers(services)

	r := inits.InitRouter(handlers)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
