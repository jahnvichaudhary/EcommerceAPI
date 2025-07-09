package internal

import (
	"context"
	"github.com/rasadov/EcommerceAPI/payment/models"
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	Close()
	GetCustomerByCustomerID(ctx context.Context, customerId string) (*models.Customer, error)
	GetCustomerByUserID(ctx context.Context, userId uint64) (*models.Customer, error)
	SaveCustomer(ctx context.Context, customer *models.Customer) error
	GetTransactionByProductID(ctx context.Context, productId string) (*models.Transaction, error)
	RegisterTransaction(ctx context.Context, transaction *models.Transaction) error
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) (Repository, error) {
	err := db.AutoMigrate(&models.Customer{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Transaction{})
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

	return &postgresRepository{db: db}, nil
}

func (repository *postgresRepository) Close() {
	sqlDB, err := repository.db.DB()
	if err == nil {
		err = sqlDB.Close()
		if err != nil {
			log.Println("Error closing postgres repository")
			log.Println(err)
		}
	}
}

func (repository *postgresRepository) GetCustomerByCustomerID(ctx context.Context, customerId string) (*models.Customer, error) {
	var customer models.Customer
	err := repository.db.WithContext(ctx).First(&customer, "id = ?", customerId).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repository *postgresRepository) GetCustomerByUserID(ctx context.Context, userId uint64) (*models.Customer, error) {
	var customer models.Customer
	err := repository.db.WithContext(ctx).First(&customer, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repository *postgresRepository) SaveCustomer(ctx context.Context, customer *models.Customer) error {
	return repository.db.WithContext(ctx).Create(&customer).Error
}

func (repository *postgresRepository) GetTransactionByProductID(ctx context.Context, productId string) (*models.Transaction, error) {
	transaction := models.Transaction{
		ProductId: productId,
	}

	err := repository.db.WithContext(ctx).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (repository *postgresRepository) RegisterTransaction(ctx context.Context, transaction *models.Transaction) error {
	return repository.db.WithContext(ctx).Create(&transaction).Error
}

func (repository *postgresRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return repository.db.WithContext(ctx).Save(&transaction).Error
}
