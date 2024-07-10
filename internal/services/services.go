package services

type AllServices struct {
	AuthService     *AuthService
	UserService     *UserService
	ItemService     *ItemService
	OrderService    *OrderService
	CategoryService *CategoryService
}
