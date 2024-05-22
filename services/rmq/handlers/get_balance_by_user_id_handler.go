package handlers

import (
	"balance-service/config"
	proto "balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	"balance-service/util"
	"fmt"
	gitProto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type GetBalanceByUserIdHandler struct {
	pgCl      *pg.PgClient
	sResponse senders.Sender
}

func NewGetBalanceByUserIdHandler(pgCl *pg.PgClient, s senders.Sender) GetBalanceByUserIdHandler {
	return GetBalanceByUserIdHandler{
		pgCl:      pgCl,
		sResponse: s,
	}
}

func (h GetBalanceByUserIdHandler) HandleMessage(body []byte) {
	req := &proto.GetBalanceByUserIdRequest{}
	err := gitProto.Unmarshal(body, req)
	log.Printf("received request: %s", req)
	util.IsError(err, "failed to unmarshal message")

	balances, err := h.pgCl.GetBalances(fmt.Sprintf(config.GetBalanceByUserIdSqlQuery, req.UserId))
	util.IsError(err, "Failed select user balance by currency and user_id")

	resp := &proto.GetBalanceByUserIdResponse{
		Id: req.Id,
		UserBalance: &proto.UserBalance{
			UserId:   req.UserId,
			Balances: balances,
		},
	}

	body, err = gitProto.Marshal(resp)
	util.IsError(err, "failed to marshal message")
	h.sResponse.SendMessage(config.RabbitBalanceExchange, config.GetBalanceByUserIdResponseRoutingKey, body)
}
