package domain

type Item struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	CategoryID    int    `json:"category_id"`
	CategoryTitle string `json:"category_title"`
}

type ItemWithAmount struct {
	Item   Item `json:"item_id"`
	Amount int  `json:"amount"`
}
