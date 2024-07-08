package services

type AllServices struct {
	UserService     *UserService
	ItemService     *ItemService
	OrderService    *OrderService
	CategoryService *CategoryService
}
