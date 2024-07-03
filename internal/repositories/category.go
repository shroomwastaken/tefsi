package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"tefsi/internal/domain"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool, all_tables *map[string]struct{}) (*CategoryRepository, error) {
	_, ok := (*all_tables)["categories"]
	if !ok {
		sqlString := `CREATE TABLE categories
		(
			id serial primary key,
			title text
		)`
		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}
	return &CategoryRepository{db: db}, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	category := &domain.Category{}
	err := r.db.QueryRow(ctx, "SELECT id, title FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Title)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	_, err := r.db.Exec(ctx, "INSERT INTO categories (title) VALUES ($1)", category.Title)
	return err
}

func (r *CategoryRepository) GetCategories(ctx context.Context) (*[]domain.Category, error) {
	var categories []domain.Category
	sqlString := "SELECT id, title FROM categories"
	rows, err := r.db.Query(ctx, sqlString)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Title)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}
	return &categories, nil
}
