package domain

type Order struct {
	ID     int    `json:"id"`
	Status Status `json:"status"`
	User   User   `json:"user"`
	Items  []Item `json:"items"`
}
