package repositories

import (
	"context"
	"fmt"

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
	item := domain.Item{}
	sql_string := `SELECT items.id, items.title, items.description, items.price, items.category, categories.title
	FROM items
	JOIN categories ON items.category = categories.id
	WHERE items.id = $1;`
	err := r.db.QueryRow(ctx, sql_string, id).Scan(
		&item.ID, &item.Title, &item.Description, &item.Price, &item.CategoryID, &item.CategoryTitle,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	sql_string := "INSERT INTO items (title, description, price, category) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(ctx, sql_string, item.Title, item.Description, item.Price, item.CategoryID)
	return err
}

func (r *ItemRepository) GetItems(ctx context.Context, filter *domain.Filter) (*[]domain.Item, error) {
	var items []domain.Item
	sql_string := `SELECT items.id, items.title, items.description, items.price, items.category, categories.title
	FROM items
	JOIN categories ON items.category = categories.id`
	if filter.CategoryID != 0 {
		sql_string += fmt.Sprintf("\nWHERE items.category = %d", filter.CategoryID)
	}
	rows, err := r.db.Query(ctx, sql_string)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := domain.Item{}
		err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Price, &item.CategoryID, &item.CategoryTitle)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &items, nil
}
