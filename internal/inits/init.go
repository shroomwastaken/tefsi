package inits

import (
	"context"
	"log"
	"tefsi/internal/handlers"
	"tefsi/internal/repositories"
	"tefsi/internal/services"

	"github.com/go-chi/chi"
)

func GetAllTables(db repositories.Pool) (map[string]struct{}, error) {
	sqlString := "SELECT table_name FROM information_schema.tables"
	rows, err := db.Query(context.Background(), sqlString)
	if err != nil {
		return nil, err
	}

	// tables := []string{}
	tables := make(map[string]struct{})

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}

		// tables = append(tables, table_name)
		tables[tableName] = struct{}{}
	}

	return tables, nil
}

func InitRepositories(db repositories.Pool, allTables map[string]struct{}) (*repositories.AllRepositories, error) {
	categoryRepo, err := repositories.NewCategoryRepository(db, &allTables)
	if err != nil {
		return nil, err
	}

	itemRepo, err := repositories.NewItemRepository(db, &allTables)
	if err != nil {
		return nil, err
	}

	userRepo, err := repositories.NewUserRepository(db, &allTables)
	if err != nil {
		log.Fatal(err)
	}

	orderRepo, err := repositories.NewOrderRepository(db, &allTables)
	if err != nil {
		return nil, err
	}

	return &repositories.AllRepositories{
		UserRepository:     userRepo,
		ItemRepository:     itemRepo,
		OrderRepository:    orderRepo,
		CategoryRepository: categoryRepo,
	}, nil
}

func InitServices(allRepos *repositories.AllRepositories) *services.AllServices {
	categoryService := services.NewDefaultCategoryService(allRepos.CategoryRepository)
	userService := services.NewDefaultUserService(allRepos.UserRepository)
	itemService := services.NewDefaultItemService(allRepos.ItemRepository)
	orderService := services.NewDefaultOrderService(allRepos.OrderRepository)

	return &services.AllServices{
		UserService:     userService,
		ItemService:     itemService,
		OrderService:    orderService,
		CategoryService: categoryService,
	}
}

func InitHandlers(allServices *services.AllServices) *handlers.AllHandlers {
	categoryHandler := handlers.NewCategoryHandler(allServices.CategoryService)
	userHandler := handlers.NewUserHandler(allServices.UserService)
	itemHandler := handlers.NewItemHandler(allServices.ItemService)
	orderHandler := handlers.NewOrderHandler(allServices.OrderService)

	return &handlers.AllHandlers{
		UserHandler:     userHandler,
		ItemHandler:     itemHandler,
		OrderHandler:    orderHandler,
		CategoryHandler: categoryHandler,
	}
}

func InitRouter(allHandlers *handlers.AllHandlers) chi.Router {
	r := chi.NewRouter()

	r.Get("/category/{id}", allHandlers.CategoryHandler.GetCategoryByID)
	r.Post("/category", allHandlers.CategoryHandler.CreateCategory)
	r.Get("/category/list", allHandlers.CategoryHandler.GetCategories)
	r.Delete("/category/delete/{id}", allHandlers.CategoryHandler.DeleteCategory)

	r.Get("/item/{id}", allHandlers.ItemHandler.GetItemByID)
	r.Post("/item", allHandlers.ItemHandler.CreateItem)
	r.Get("/item/list", allHandlers.ItemHandler.GetItems)
	r.Delete("/item/delete/{id}", allHandlers.ItemHandler.DeleteItem)

	r.Get("/users/{id}", allHandlers.UserHandler.UserRequired(allHandlers.UserHandler.GetUserByID))
	r.Post("/users", allHandlers.UserHandler.CreateUser)
	r.Post("/users/login", allHandlers.UserHandler.Login)
	r.Delete("/users/delete/{id}", allHandlers.UserHandler.DeleteUser)

	r.Get("/order/{id}", allHandlers.OrderHandler.GetOrderByID)
	r.Post("/order", allHandlers.OrderHandler.CreateOrder)
	r.Get("/order/list", allHandlers.UserHandler.AdminRequired(allHandlers.OrderHandler.GetOrders))
	r.Get("/order/list/{id}", allHandlers.UserHandler.UserRequired(allHandlers.OrderHandler.GetOrdersByUserID))
	r.Delete("/order/delete/{id}", allHandlers.OrderHandler.DeleteOrder)

	return r
}
