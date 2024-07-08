package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"tefsi/internal/domain"
)

// Реализация репозитория
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool, allTables *map[string]struct{}) (*UserRepository, error) {
	_, ok := (*allTables)["users"]

	if !ok {
		sqlString := `CREATE TABLE users
        (
            id serial primary key,
            name text,
            email text
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
            user_id int,
            FOREIGN KEY item REFERENCES items(id),
            FOREIGN KEY user_id REFERENCES users(id),
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
	err := r.db.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	return err
}

func (r *UserRepository) GetUserCartByID(ctx context.Context, id int) (*[]domain.Item, error) {
	cart := make([]domain.Item, 0)
	return &cart, nil // plug
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
