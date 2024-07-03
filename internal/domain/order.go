package domain

// type Order struct {
// 	ID          int           `json:"id"`
// 	StatusID    int           `json:"status_id"`
// 	StatusTitle string        `json:"status_title"`
// 	UserID      int           `json:"user_id"`
// 	Items       map[*Item]int `json:"items"`
// }

type Order struct {
	ID          int    `json:"id"`
	StatusID    int    `json:"status_id"`
	StatusTitle string `json:"status_title"`
	UserID      int    `json:"user_id"`
	Items       []Item `json:"items"`
	Amounts     []int  `json:"amounts"`
}
