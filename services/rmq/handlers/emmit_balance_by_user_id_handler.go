package handlers

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	"balance-service/util"
	"fmt"
	gitProto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

type EmmitBalanceByUserIdHandler struct {
	pgCl      *pg.PgClient
	sResponse senders.Sender
}

func NewEmmitBalanceByUserIdHandler(pgCl *pg.PgClient, s senders.Sender) EmmitBalanceByUserIdHandler {
	return EmmitBalanceByUserIdHandler{
		pgCl:      pgCl,
		sResponse: s,
	}
}

func (h EmmitBalanceByUserIdHandler) HandleMessage(body []byte) {
	req := &proto.EmmitBalanceByUserIdRequest{}

	err := gitProto.Unmarshal(body, req)
	util.IsError(err, "Failed to unmarshal message")

	h.pgCl.Exec(config.EmitBalanceByUserIdSqlQuery,
		req.Currency,
		fmt.Sprintf("+%g", req.Amount),
		0,
		time.Now().Unix(),
		req.UserId,
	)

	balances, err := h.pgCl.GetBalances(fmt.Sprintf(config.GetBalanceByUserIdCurrencySqlQuery, req.UserId, req.Currency))
	util.IsError(err, "Failed select user balance by currency and user_id")

	resp := &proto.UserBalance{
		UserId:   req.UserId,
		Balances: balances,
	}
	body, err = gitProto.Marshal(resp)
	util.IsError(err, "Failed to marshal message")

	h.sResponse.SendMessage("e.balances.forward", "r.balance.EmitUserBalanceResponse", body)
	log.Printf("Received OrderUpdateEvent: %+v\n", req)
}
