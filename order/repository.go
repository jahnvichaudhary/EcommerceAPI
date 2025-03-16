package order

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(databaseURl string) (Repository, error) {
	db, err := sql.Open("postgres", databaseURl)
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r postgresRepository) Close() {
	r.db.Close()
}

func (r postgresRepository) PutOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Here we insert general info about the order: created at, account id, total price
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		order.ID,
		order.CreatedAt,
		order.AccountID,
		order.TotalPrice,
	)
	if err != nil {
		return err
	}

	// Insert order products. We iterate over products in order and add them to table order_products
	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, product := range order.Products {
		_, err = stmt.ExecContext(ctx, order.ID, product.ID, product.Quantity)
		if err != nil {
			return err
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return stmt.Close()
}

func (r postgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
		o.id,
		o.created_at,
		o.account_id,
		o.total_price::money::numeric::float8,
		op.product_id,
		op.quantity
		FROM orders o JOIN order_products op ON (o.id = op.order_id)
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accountId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	var products []OrderedProduct
	var lastOrderID string
	order := &Order{}
	orderedProduct := &OrderedProduct{}

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrderID != "" && lastOrderID != order.ID {
			order.Products = products
			orders = append(orders, *order)
			products = []OrderedProduct{}
		}

		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		lastOrderID = order.ID
	}

	if lastOrderID != "" {
		order.Products = products
		orders = append(orders, *order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
