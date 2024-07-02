package repositories

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"tefsi/internal/domain"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetItemByID(ctx context.Context, id int) (*domain.Item, error) {
	item := &domain.Item{}
	sql_string := "SELECT items.id, items.title, items.description, items.price, category.id, category.title FROM items WHERE id = $1 INNER JOIN categories ON items.category = categories.id"
	err := r.db.QueryRow(ctx, sql_string, id).Scan(
		&item.ID, &item.Title, &item.Description, &item.Price, &item.Category.ID, &item.Category.Title,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	sql_string := "INSERT INTO item (title, description, price, category) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(ctx, sql_string, item.Title, item.Description, item.Price, item.Category.ID)
	return err
}
