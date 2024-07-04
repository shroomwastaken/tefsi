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
	// TODO: needs sooooo many checks
	sqlString := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, sqlString, id)
	return err
}
