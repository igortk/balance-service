package pg

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/services/rmq/handlers"
	"balance-service/util"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(cfg *config.PostreSqlConfig) (*Client, error) {
	db, err := sqlx.Connect(
		config.DriverName,
		fmt.Sprintf(config.PgConnectionUrlPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))

	if err != nil {
		return nil, fmt.Errorf("failed connect to postgre: %v", err)
	}

	return &Client{
		db: db,
	}, nil
}

func (cl *Client) EmitCurrency(userId, currencyName string, amount float64) error {
	tx, err := cl.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	result, err := tx.Exec(EmitBalanceByUserIdSqlQuery, currencyName, amount, 0, time.Now().Unix(), userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute balance update: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot get rows affected: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (cl *Client) GetUserBalance(userId, currencyName string) (*proto.Balance, error) {
	balance := make([]proto.Balance, 0, 1)

	err := cl.db.Select(&balance, GetBalanceByUserIdCurrencySqlQuery, userId, currencyName)
	if err != nil {
		return nil, fmt.Errorf("failed get user balance: %w", err)
	}

	return &balance[0], nil
}

func (cl *Client) GetUserBalances(userId string) ([]*proto.Balance, error) {
	balances := make([]*proto.Balance, 0, 1)

	err := cl.db.Select(&balances, GetBalanceByUserIdSqlQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("failed get user balances: %w", err)
	}

	return balances, nil
}

func (cl *Client) UpdateBalancesTx(users ...*handlers.User) (err error) {
	ctx := context.Background()

	tx, err := cl.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	now := time.Now().Unix()
	for _, u := range users {
		if u.CurrencyName != "" && u.Balance != 0 {
			_, err = tx.ExecContext(ctx, UpdateBalanceByUserIdSqlQuery, u.CurrencyName, u.Balance, 0, now, u.UserId)
			if err != nil {
				return fmt.Errorf("balance update failed: %w", err)
			}
		}

		if u.LockedCurrencyName != "" && u.LockedBalance != 0 {
			_, err = tx.ExecContext(ctx, UpdateBalanceByUserIdSqlQuery, u.LockedCurrencyName, 0, u.LockedBalance, now, u.UserId)
			if err != nil {
				return fmt.Errorf("locked balance update failed: %w", err)
			}
		}
	}

	return nil
}

func (cl *Client) Exec(query string, args ...interface{}) interface{} {
	result, err := cl.db.Exec(query, args...)
	util.IsError(err, "Failed insert")
	return result
}
