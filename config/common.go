package config

// RMQ exchange
const (
	RabbitEventsExchange  = "e.events.forward"
	RabbitBalanceExchange = "e.balances.forward"
)

// RMQ rk
const (
	GetBalanceByUserIdRequestRoutingKey  = "r.balance-service.balances.#.GetBalanceByUserIdRequest"
	EmitBalanceByUserIdRequestRoutingKey = "r.balance-service.balances.#.EmmitUserBalanceRequest"
	UpdatedOrderEventRoutingKey          = "r.event.order.OrderUpdateEvent"
	GetBalanceByUserIdResponseRoutingKey = "r.balance.GetBalanceByUserIdResponse"
)

// Queue name
const (
	UpdatedOrderEventQueueName         = "q.balance-service.order.event"
	GetBalanceByUserIdRequestQueueName = "q.balance-service.user.balance.get.request"
	EmitUserBalanceRequestQueueName    = "q.balance-service.user.balance.emit.request"
)

// Errors
const (
	ErrLoadConfig = "Error load configuration"
	ErrParseLog   = "Error parse log level"
	ErrConnectDb  = "Error connect db"
)

const (
	//GetBalanceByUserIdSqlQuery         = "SELECT * FROM balances b WHERE b.user_id = %s"
	GetBalanceByUserIdCurrencySqlQuery = "SELECT b.currency_id as currency,\nb.balance as balance,\nb.locked_balance as locked_balance,\nb.updated_date as updated_date\nFROM balances b WHERE b.user_id = '%s' and b.currency_id = '%s'"
	GetBalanceByUserIdSqlQuery         = "SELECT b.currency_id as currency,\nb.balance as balance,\nb.locked_balance as locked_balance,\nb.updated_date as updated_date\nFROM balances b WHERE b.user_id = '%s'"
	EmmitBalanceByUserIdSqlQuery       = "INSERT INTO balances (currency_id, balance, locked_balance, updated_date, user_id) VALUES ($1, $2, $3, $4, $5)\nON CONFLICT (currency_id, user_id) DO UPDATE\nSET balance = balances.balance + $2, locked_balance = balances.locked_balance + $3;"
)

// Url pattern
const (
	PgConnectionUrlPattern  = "postgresql://%s:%s@%s:%d/%s"
	RmqUrlConnectionPattern = "amqp://%s:%s@%s:%d/"
)

const (
	DriverName = "pgx"
)
