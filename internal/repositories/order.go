package repositories

import (
	"context"
	"tefsi/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool, all_tables *map[string]struct{}) (*OrderRepository, error) {
	_, ok := (*all_tables)["statuses"]
	if !ok {
		sqlString := `CREATE TABLE statuses
            (
                id serial primary key,
                title text
            )`
		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}

	_, ok = (*all_tables)["orders"]
	if !ok {
		sqlString := `CREATE TABLE orders
        (
            id serial primary key,
            status int,
            user_id int,
            FOREIGN KEY (status) REFERENCES statuses(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`
		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}

	_, ok = (*all_tables)["items_orders"]
	if !ok {
		sqlString := `CREATE TABLE items_orders
        (
            id serial primary key,
            item int,
            amount int,
            order_id int
        )`

		_, err := db.Exec(context.Background(), sqlString)
		if err != nil {
			return nil, err
		}
	}

	return &OrderRepository{db: db}, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	orderSQL := "INSERT INTO orders (id, status, user_id) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, orderSQL, order.ID, order.StatusID, order.UserID)
	if err != nil {
		return err
	}

	itemSQL := "INSERT into items_orders (item, order_id, amount) VALUES ($1, $2, $3)"

	for i := range order.Items {
		_, err := r.db.Exec(ctx, itemSQL, order.Items[i].ID, order.ID, order.Amounts[i])
		if err != nil {
			return err
		}
	}

	return err
}

// s dnem prikolov
func (r *OrderRepository) getStatusTitleAndItems(ctx context.Context, order *domain.Order) (string, *[]domain.Item, *[]int, error) {
	var statusTitle string
	// items := make(map[*domain.Item]int)
	items := []domain.Item{}
	amounts := []int{}

	statusTitleSQL := `SELECT statuses.title
    FROM statuses
    WHERE id = $1`

	err := r.db.QueryRow(ctx, statusTitleSQL, order.StatusID).Scan(&statusTitle)
	if err != nil {
		return "", nil, nil, err
	}

	itemsSQL := `SELECT items.id, items.title, items.description, items.price, items.category, categories.title, items_orders.amount
    FROM items_orders
    JOIN items ON items.id = items_orders.item
    JOIN categories ON items.category = categories.id
    WHERE items_orders.order_id = $1`

	itemsRows, err := r.db.Query(ctx, itemsSQL, order.ID)

	// TODO: proper error handling
	if err != nil {
		return "", nil, nil, err
	}

	for itemsRows.Next() {
		item := domain.Item{}
		var amount int

		err := itemsRows.Scan(&item.ID, &item.Title, &item.Description, &item.Price, &item.CategoryID, &item.CategoryTitle, &amount)

		if err != nil {
			return "", nil, nil, err
		}

		items = append(items, item)
		amounts = append(amounts, amount)
	}

	return statusTitle, &items, &amounts, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id int) (*domain.Order, error) {
	order := domain.Order{}

	sql_string := `SELECT orders.id, orders.status, orders.user_id
    FROM orders
    WHERE orders.id = $1`

	err := r.db.QueryRow(ctx, sql_string, id).Scan(&order.ID, &order.StatusID, &order.UserID)
	if err != nil {
		return nil, err
	}

	statusTitle, items, amounts, err := r.getStatusTitleAndItems(ctx, &order)

	if err != nil {
		return nil, err
	}

	order.StatusTitle = statusTitle
	order.Items = *items
	order.Amounts = *amounts

	return &order, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context) (*[]domain.Order, error) {
	var orders []domain.Order

	sqlString := "Select orders.id, orders.status, orders.user_id FROM orders"

	rows, err := r.db.Query(ctx, sqlString)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(&order.ID, &order.StatusID, &order.UserID)
		if err != nil {
			return nil, err
		}

		statusTitle, items, amounts, err := r.getStatusTitleAndItems(ctx, &order)
		if err != nil {
			return nil, err
		}

		order.StatusTitle = statusTitle
		order.Items = *items
		order.Amounts = *amounts

		orders = append(orders, order)
	}

	return &orders, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	ordersSQL := `UPDATE orders
    SET orders.status = $1, orders.user_id = $2
    WHERE orders.id = $3`

	_, err := r.db.Exec(ctx, ordersSQL, order.StatusID, order.StatusTitle, order.ID)
	if err != nil {
		return err
	}

	deleteItemsSQL := `DELETE FROM items_orders
    WHERE order_id = $1`

	_, err = r.db.Exec(ctx, deleteItemsSQL, order.ID)
	if err != nil {
		return err
	}

	addItemSQL := "INSERT into items_orders (item, order_id, amount) VALUES ($1, $2, $3)"

	// for item, amount := range order.Items {
	// 	_, err := r.db.Exec(ctx, addItemSQL, item.ID, order.ID, amount)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	for i := range order.Items {
		_, err := r.db.Exec(ctx, addItemSQL, order.Items[i].ID, order.ID, order.Amounts[i])
		if err != nil {
			return err
		}
	}

	return nil
}
