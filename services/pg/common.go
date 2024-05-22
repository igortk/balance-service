package pg

// Queries

const (
	InsertUserBalanceSQLQueryPattern = "INSERT INTO public.balances(currency, balance, locked_balance, updated_date, user_id)VALUES($1, $2, $3, $4, $5)"
	UpdateBalanceByUserIdSqlQuery    = "INSERT INTO balances (currency_id, balance, locked_balance, updated_date, user_id) VALUES ($1, $2, $3,$4,$5)\nON CONFLICT (currency_id, user_id) DO UPDATE\nSET balance = balances.balance + $2, locked_balance = balances.locked_balance + $3;"
)
