package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	"hse_school/internal/handlers"
	"hse_school/internal/repositories"
	"hse_school/internal/services"
)

func main() {
	postgresUser := "postgres"
	postgresPassword := "password"
	postgresDB := "postgres"    // Имя базы данных, которое вы хотите использовать
	postgresHost := "localhost" // Хост базы данных
	postgresPort := "5432"      // Порт базы данных

	// Формирование строки соединения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	// Подключение к базе данных
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание репозитория, сервиса и обработчиков
	repo := repositories.NewUserRepository(db)
	service := services.NewDefaultUserService(repo)
	handler := handlers.NewUserHandler(service)

	// Настройка маршрутизатора
	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserByID)
	r.Post("/users", handler.CreateUser)

	// Запуск HTTP сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
