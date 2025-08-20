package config

// RMQ exchange
const (
	RabbitEventsExchange  = "e.events.forward"
	RabbitBalanceExchange = "e.balances.forward"
)

// RMQ rk
const (
	GetBalanceByUserIdRequestRoutingKey  = "r.balance-service.balances.#.GetBalanceByUserIdRequest"
	EmitBalanceByUserIdRequestRoutingKey = "r.balance-service.balances.#.EmitUserBalanceRequest"
	UpdatedOrderEventRoutingKey          = "r.ops.balance-service.order.OrderUpdateEvent"
	GetBalanceByUserIdResponseRoutingKey = "r.balance.GetBalanceByUserIdResponse"
)

// queue name
const (
	UpdatedOrderEventQueueName         = "q.balance-service.order.update-event"
	GetBalanceByUserIdRequestQueueName = "q.balance-service.user.balance.get.request"
	EmitUserBalanceRequestQueueName    = "q.balance-service.user.balance.emit.request"
)

// Errors
const (
	ErrLoadConfig = "Error load configuration"
	ErrParseLog   = "Error parse log level"
	ErrConnectDb  = "Error connect db"
	ErrConnectRmq = "Error connect RabbitMq"
)

// Url pattern
const (
	PgConnectionUrlPattern  = "postgresql://%s:%s@%s:%d/%s"
	RmqUrlConnectionPattern = "amqp://%s:%s@%s:%d/"
)

const (
	DriverName = "pgx"
)
