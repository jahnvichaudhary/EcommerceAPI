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
	GetCustomerByUserID(ctx context.Context, userId int) (*models.Customer, error)
	RegisterTransaction(ctx context.Context, transaction *models.Transaction) error
	SaveCustomer(ctx context.Context, customer *models.Customer) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
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

func (repository *postgresRepository) GetCustomerByUserID(ctx context.Context, userId int) (*models.Customer, error) {
	var customer models.Customer
	err := repository.db.WithContext(ctx).First(&customer, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repository *postgresRepository) RegisterTransaction(ctx context.Context, transaction *models.Transaction) error {
	return repository.db.WithContext(ctx).Create(&transaction).Error
}

func (repository *postgresRepository) SaveCustomer(ctx context.Context, customer *models.Customer) error {
	return repository.db.WithContext(ctx).Create(&customer).Error
}
