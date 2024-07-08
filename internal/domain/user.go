package domain

// Структуры данных
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
