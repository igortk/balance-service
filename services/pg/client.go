package pg

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/util"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"reflect"
	"time"
)

type PgClient struct {
	db *sqlx.DB
}

func NewClient(cfg *config.PostreSqlConfig) *PgClient {
	db, err := sqlx.Connect(
		config.DriverName,
		fmt.Sprintf(config.PgConnectionUrlPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))
	util.IsError(err, "err pg connect")

	return &PgClient{
		db: db,
	}
}

func (cl *PgClient) EmitCurrency(userId, currencyName string, amount float64) error {
	tx, err := cl.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	result, err := tx.Exec(config.EmitBalanceByUserIdSqlQuery, currencyName, amount, 0, time.Now().Unix(), userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute balance update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot get rows affected: %w", err)
	}

	if rows != 1 {
		tx.Rollback()
		return fmt.Errorf("expected 1 row to be affected, but got %d", rows)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (cl *PgClient) GetUserBalances(userId, currencyName string) (*proto.Balance, error) {
	balance := &proto.Balance{}

	err := cl.db.Select(balance, config.GetBalanceByUserIdCurrencySqlQuery, userId, currencyName)
	if err != nil {
		return nil, fmt.Errorf("failed get user balance: %w", err)
	}

	return balance, nil
}

func (cl *PgClient) Exec(query string, args ...interface{}) interface{} {
	result, err := cl.db.Exec(query, args...)
	util.IsError(err, "Failed insert")
	return result
}

func (cl *PgClient) Select(query string, resp interface{}, args ...interface{}) interface{} {
	result := cl.db.Select(resp, query)
	return result
}
func (cl *PgClient) GetBalances(query string) ([]*proto.Balance, error) {
	rows, err := cl.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*proto.Balance
	for rows.Next() {
		balance := &proto.Balance{}
		err := rows.Scan(&balance.Currency, &balance.Balance, &balance.LockedBalance, &balance.UpdatedDate)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}
func (cl *PgClient) Query(query string, dest interface{}, args ...interface{}) error {
	rows, err := cl.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	destValue := reflect.ValueOf(dest)
	destType := destValue.Type()
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(interface{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		// Создаем новый экземпляр структуры
		newStruct := reflect.New(destType.Elem())

		// Заполняем поля структуры значениями из запроса
		for i, column := range columns {
			field := newStruct.Elem().FieldByName(column)
			if !field.IsValid() {
				return fmt.Errorf("struct field %s not found", column)
			}
			val := *values[i].(*interface{})
			if val != nil {
				field.Set(reflect.ValueOf(val).Elem())
			}
		}

		// Добавляем структуру к результатам
		destValue.Elem().Set(reflect.Append(destValue.Elem(), newStruct.Elem()))
	}

	return nil
}

/*
func (cl *PgClient) GetBalanceByUserId(req *proto.GetBalanceByUserIdRequest) ([]*proto.Balance, error) {
	//cl.db.Select(, buildQueryByProto(req)
	return nil, nil
}*/
/*
func BuildQueryByProto(req *proto.GetBalanceByUserIdRequest) string {
	return fmt.Sprintf(config.GetBalanceByUserIdSqlQuery, req.GetUserId())
}*/
