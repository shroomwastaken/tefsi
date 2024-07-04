package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

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
			login text unique,
			password text,
			is_admin bool
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
	_, err = r.db.Exec(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", user.Login, user.Password)
	return err
}

func (r *UserRepository) GetUserCartByID(ctx context.Context, id int) (*[]domain.Item, error) {
	cart := make([]domain.Item, 0)
	return &cart, nil // plug
}

func (r *UserRepository) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *UserRepository) CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
