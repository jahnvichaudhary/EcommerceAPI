package order

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(databaseURl string) (Repository, error) {
	db, err := gorm.Open(postgres.Open(databaseURl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Order{}, &ProductsInfo{})

	return &postgresRepository{db}, nil
}

func (repository *postgresRepository) Close() {
	sqlDB, err := repository.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (repository *postgresRepository) PutOrder(ctx context.Context, order Order) error {
	tx := repository.db.WithContext(ctx).Begin()

	err := tx.WithContext(ctx).Create(&order).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, product := range order.Products {
		orderedProduct := ProductsInfo{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  int(product.Quantity),
		}
		err = tx.Create(&orderedProduct).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (repository *postgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	var orders []Order
	err := repository.db.WithContext(ctx).
		Table("orders o").
		Select("o.id, o.created_at, o.account_id, o.total_price::money::numeric::float8, op.product_id, op.quantity").
		Joins("JOIN order_products op on o.id = op.order_id").
		Where("o.account_id = ?", accountId).
		Order("o.id").
		Scan(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}
