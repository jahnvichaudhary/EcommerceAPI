package account

import (
	"context"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{r}
}

func (service accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	acc := Account{
		Name: name,
	}
	err := service.repository.PutAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (service accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return service.GetAccount(ctx, id)
}

func (service accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return service.GetAccounts(ctx, skip, take)

}
