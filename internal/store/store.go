package store

import (
	"github.com/Khvan-Group/common-library/errors"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"wallet-service/internal/clients"
	users "wallet-service/internal/common/model"
	"wallet-service/internal/database"
	"wallet-service/internal/models"
)

type WalletStore struct {
	db     *sqlx.DB
	client *resty.Client
}

func New() *WalletStore {
	return &WalletStore{
		db:     database.DB,
		client: resty.New(),
	}
}

func (s *WalletStore) Save(input models.Wallet) *errors.CustomError {
	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var existsWallet bool
		tx.Get(&existsWallet, "select exists(select 1 from t_wallets where username = $1)", input.Username)

		if !clients.ExistsUser(input.Username, s.client) {
			return errors.NewBadRequest("Пользователь не существует.")
		}

		if existsWallet {
			_, err := tx.NamedExec("update t_wallets set total = :total where username = :username", input)
			if err != nil {
				return errors.NewInternal("Failed to update wallet transaction")
			}

			return nil
		}

		_, err := tx.NamedExec("insert into t_wallets(username, total) values (:username, :total)", input)
		if err != nil {
			return errors.NewInternal("Failed to insert wallet transaction")
		}

		return nil
	})
}

func (s *WalletStore) Update(input models.WalletUpdate) *errors.CustomError {
	return database.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var entity models.Wallet
		err := tx.Get(&entity, "select * from t_wallets where username = $1", input.Username)
		if err != nil {
			return errors.NewBadRequest("Пользователь не имеет кошелька.")
		}

		if input.Action == models.WALLET_TOTAL_ADD {
			input.Total += entity.Total
		} else if input.Action == models.WALLET_TOTAL_SUBSTRUCT {
			if input.Total > entity.Total {
				return errors.NewBadRequest("Невозможно списать сумму, превышающую имеющуюся.")
			}

			input.Total = entity.Total - input.Total
		} else {
			return errors.NewBadRequest("Неверное действие над кошельком.")
		}

		_, err = tx.NamedExec("update t_wallets set total = :total where username = :username", input)
		if err != nil {
			return errors.NewInternal("Failed to insert wallet transaction")
		}

		return nil
	})
}

func (s *WalletStore) FindByUser(user users.JwtUser) (*models.Wallet, *errors.CustomError) {
	var wallet models.Wallet
	err := s.db.Get(&wallet, "select * from t_wallets where username = $1", user.Login)
	if err != nil {
		return nil, errors.NewBadRequest("У пользователя нет кошелька.")
	}

	return &wallet, nil
}

func (s *WalletStore) Delete(user users.JwtUser) *errors.CustomError {
	_, err := s.db.Exec("delete from t_wallets where username = $1", user.Login)
	if err != nil {
		return errors.NewInternal("Failed to delete wallet transaction")
	}

	return nil
}
