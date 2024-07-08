package handlers

type AllHandlers struct {
	UserHandler     *UserHandler
	ItemHandler     *ItemHandler
	OrderHandler    *OrderHandler
	CategoryHandler *CategoryHandler
}
