package repositories

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"tefsi/internal/domain"
)

type UserRepository struct {
	db Pool
}

func NewUserRepository(db Pool, allTables *map[string]struct{}) (*UserRepository, error) {
	_, ok := (*allTables)["users"]
	if !ok {
		sqlString := `CREATE TABLE users
		(
			id serial primary key,
			login text unique,
			password text,
			is_admin bool
		)`
		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}

	_, ok = (*allTables)["items_users"]
	if !ok {
		sqlString := `CREATE TABLE items_users
        (
            id serial primary key,
            item int,
            amount int,
            user_id int,
            FOREIGN KEY (item) REFERENCES items(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`
		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(ctx, "SELECT id, login, password, is_admin FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Login, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UserExists(ctx context.Context, login string) error {
	rows, err := r.db.Query(ctx, "SELECT login FROM users WHERE login = $1", login)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return fmt.Errorf("user does not exist")
	}
	return nil
}

func (r *UserRepository) CheckUserByDomain(ctx context.Context, user *domain.User) error {
	var correctPassword string
	err := r.db.QueryRow(ctx, "SELECT password FROM users WHERE login = $1", user.Login).
		Scan(&correctPassword)
	if err != nil {
		return err
	}
	err = r.CheckPasswordHash((*user).Password, correctPassword)
	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	var err error
	user.Password, err = r.HashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "INSERT INTO users (login, password, is_admin) VALUES ($1, $2, $3)", user.Login, user.Password, user.IsAdmin)
	return err
}

// TODO: check that it works
func (r *UserRepository) GetUserCartByID(ctx context.Context, id int) (*[]domain.ItemWithAmount, error) {
	sqlString := `SELECT items.id, items.title, items.description, items.price, items.category, categories.title, items_users.amount
    FROM items_users
    JOIN items ON items.id = items_users.id
    JOIN categories on items.category = categories.id
    WHERE items_users.user = $1`

	rows, err := r.db.Query(ctx, sqlString, id)
	if err != nil {
		return nil, err
	}

	items := []domain.ItemWithAmount{}

	for rows.Next() {
		item := domain.ItemWithAmount{}

		err := rows.Scan(&item.Item.ID, &item.Item.Title, &item.Item.Description, &item.Item.Price, &item.Item.CategoryID, &item.Item.CategoryTitle, &item.Amount)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &items, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	deleteUserSQL := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, deleteUserSQL, id)
	if err != nil {
		return err
	}

	deleteItemsUsersSQL := "DELETE FROM items_users WHERE user = $1"
	_, err = r.db.Exec(ctx, deleteItemsUsersSQL, id)
	if err != nil {
		return err
	}

	selectOrdersSQL := "SELECT id FROM orders WHERE user = $1"
	ordersRows, err := r.db.Query(ctx, selectOrdersSQL, id)
	orderIDs := []int{}

	for ordersRows.Next() {
		var orderID int
		err := ordersRows.Scan(&orderID)
		if err != nil {
			return err
		}

		orderIDs = append(orderIDs, orderID)
	}

	deleteOrdersSQL := "DELETE FROM orders WHERE user = $1"
	_, err = r.db.Exec(ctx, deleteOrdersSQL, id)
	if err != nil {
		return err
	}

	deleteItemsOrdersSQL := "DELETE FROM items_orders WHERE order = $1"
	for _, orderID := range orderIDs {
		_, err := r.db.Exec(ctx, deleteItemsOrdersSQL, orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) UserIsAdmin(ctx context.Context, login string) (bool, error) {
	var isAdmin bool
	err := r.db.QueryRow(ctx, "SELECT is_admin FROM users WHERE login = %1", login).Scan(isAdmin)
	return isAdmin, err
}

func (r *UserRepository) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *UserRepository) CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
