package pg

// Queries

const (
	InsertUserBalanceSQLQueryPattern   = "INSERT INTO public.balances(currency, balance, locked_balance, updated_date, user_id)VALUES($1, $2, $3, $4, $5)"
	UpdateBalanceByUserIdSqlQuery      = "INSERT INTO balances (currency_id, balance, locked_balance, updated_date, user_id) VALUES ($1, $2, $3,$4,$5)\nON CONFLICT (currency_id, user_id) DO UPDATE\nSET balance = balances.balance + $2, locked_balance = balances.locked_balance + $3;"
	GetBalanceByUserIdCurrencySqlQuery = "SELECT c2.name as currency,\nb.balance as balance,\nb.locked_balance as lockedBalance,\nb.updated_date as updatedDate\nFROM balances b join currencies c2 on b.currency_id = c2.id WHERE b.user_id = $1 and b.currency_id = (select c.id from currencies c where c.name = $2)"
	GetBalanceByUserIdSqlQuery         = "SELECT c2.name as currency,\nb.balance as balance,\nb.locked_balance as lockedBalance,\nb.updated_date as updatedDate\nFROM balances b join currencies c2 on b.currency_id = c2.id WHERE b.user_id = $1"
	EmitBalanceByUserIdSqlQuery        = "INSERT INTO balances (currency_id, balance, locked_balance, updated_date, user_id)\nSELECT c.id, $2, $3, $4, $5\nFROM currencies c\nWHERE c.name = $1\nON CONFLICT (currency_id, user_id) DO UPDATE\nSET balance = balances.balance + EXCLUDED.balance,\n    locked_balance = balances.locked_balance + EXCLUDED.locked_balance;"
)
