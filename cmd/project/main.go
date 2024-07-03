package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	"tefsi/internal/handlers"
	"tefsi/internal/repositories"
	"tefsi/internal/services"
)

func GetAllTables(db *pgxpool.Pool) (map[string]struct{}, error) {
	sqlString := "SELECT table_name FROM information_schema.tables"
	rows, err := db.Query(context.Background(), sqlString)
	if err != nil {
		return nil, err
	}

	// tables := []string{}
	tables := make(map[string]struct{})

	for rows.Next() {
		var table_name string
		err := rows.Scan(&table_name)
		if err != nil {
			return nil, err
		}

		// tables = append(tables, table_name)
		tables[table_name] = struct{}{}
	}

	return tables, nil
}

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

	all_tables, err := GetAllTables(db)
	if err != nil {
		panic(err)
	}

	// Создание репозитория, сервиса и обработчиков
	repo := repositories.NewUserRepository(db, &all_tables)
	service := services.NewDefaultUserService(repo)
	handler := handlers.NewUserHandler(service)

	// Настройка маршрутизатора
	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserByID)
	r.Post("/users", handler.CreateUser)

	categoryRepo := repositories.NewCategoryRepository(db, &all_tables)
	categoryService := services.NewDefaultCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	r.Get("/category/{id}", categoryHandler.GetCategoryByID)
	r.Post("/category", categoryHandler.CreateCategory)
	r.Get("/category/list", categoryHandler.GetCategories)

	itemRepo := repositories.NewItemRepository(db, &all_tables)
	itemService := services.NewDefaultItemService(itemRepo)
	itemHandler := handlers.NewItemHandler(itemService)

	r.Get("/item/{id}", itemHandler.GetItemByID)
	r.Post("/item", itemHandler.CreateItem)
	r.Get("/item/list", itemHandler.GetItems)

	orderRepo := repositories.NewOrderRepository(db, &all_tables)
	orderService := services.NewDefaultOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	r.Get("/order/{id}", orderHandler.GetOrderByID)
	r.Post("/order", orderHandler.CreateOrder)
	r.Get("/order/list", orderHandler.GetOrders)

	// Запуск HTTP сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
