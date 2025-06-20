package pg

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/util"
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
	/*
		if rows != 1 {
			tx.Rollback()
			return fmt.Errorf("expected 1 row to be affected, but got %d", rows)
		}*/

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
	var balances []*proto.Balance

	err := cl.db.Select(balances, GetBalanceByUserIdSqlQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("failed get user balances: %w", err)
	}

	return balances, nil
}

func (cl *Client) Exec(query string, args ...interface{}) interface{} {
	result, err := cl.db.Exec(query, args...)
	util.IsError(err, "Failed insert")
	return result
}
