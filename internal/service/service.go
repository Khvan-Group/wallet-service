package service

import (
	"github.com/Khvan-Group/common-library/errors"
	users "wallet-service/internal/common/model"
	"wallet-service/internal/models"
	"wallet-service/internal/store"
)

type WalletService interface {
	Save(input models.Wallet) *errors.CustomError
	Update(input models.WalletUpdate) *errors.CustomError
	FindByUser(user users.JwtUser) (*models.Wallet, *errors.CustomError)
	Delete(user users.JwtUser) *errors.CustomError
}

type Wallets struct {
	Service WalletService
}

func New() *Wallets {
	return &Wallets{
		Service: store.New(),
	}
}
