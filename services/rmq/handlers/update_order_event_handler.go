package handlers

import (
	"balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/util"
	gitProto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type UpdateOrderEventHandler struct {
	pgCl *pg.PgClient
}

func NewUpdateOrderEventHandler(pgCl *pg.PgClient) *UpdateOrderEventHandler {
	return &UpdateOrderEventHandler{
		pgCl: pgCl,
	}
}

func (h UpdateOrderEventHandler) HandleMessage(body []byte) {
	event := &proto.OrderUpdateEvent{}
	err := gitProto.Unmarshal(body, event)
	util.IsError(err, "Failed to unmarshal message")
	log.Printf("received event: %s", event)

	if event.Error != nil {
		return
	}

	currency := ""

	switch event.Order.Direction {
	case proto.Direction_ORDER_DIRECTION_BUY:
		currency = strings.Split(event.Order.Pair, "/")[1]
		break
	case proto.Direction_ORDER_DIRECTION_SELL:
		currency = strings.Split(event.Order.Pair, "/")[0]
		break
	default:
		log.Errorf("unsupported direction value")
		return
	}

	if event.Order.Status == proto.OrderStatus_ORDER_STATUS_REMOVED {
		h.pgCl.Exec(
			pg.UpdateBalanceByUserIdSqlQuery,
			currency,
			event.Order.InitPrice*(event.Order.InitVolume-event.Order.FillVolume),
			(-1)*event.Order.InitPrice*(event.Order.InitVolume-event.Order.FillVolume),
			time.Now().Unix(),
			event.Order.UserId,
		)
		log.Info("successful update of user balance")
		return
	}

	h.pgCl.Exec(
		pg.UpdateBalanceByUserIdSqlQuery,
		currency,
		(-1)*event.Order.InitPrice*event.Order.InitVolume,
		event.Order.InitPrice*event.Order.InitVolume,
		time.Now().Unix(),
		event.Order.UserId,
	)
	log.Info("successful update of user balance")
}
