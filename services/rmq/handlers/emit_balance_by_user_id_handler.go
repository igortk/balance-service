package handlers

import (
	"balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	log "github.com/sirupsen/logrus"
)

type EmitBalanceByUserIdHandler struct {
	pgCl      *pg.PgClient
	sResponse *senders.Sender
}

func NewEmitBalanceByUserIdHandler(pgCl *pg.PgClient, s *senders.Sender) *EmitBalanceByUserIdHandler {
	return &EmitBalanceByUserIdHandler{
		pgCl:      pgCl,
		sResponse: s,
	}
}

func (h *EmitBalanceByUserIdHandler) HandleMessage(body []byte) {
	req, err := unmarshalEmitBalanceByUserIdRequest(&body)
	if err != nil {
		log.Errorf("Failed deserialize request: %v", err)
		return
	}

	//TODO add validate request

	resp := &proto.EmitBalanceByUserIdResponse{Id: req.Id, UserId: req.UserId}

	err = h.pgCl.EmitCurrency(req.UserId, req.CurrencyName, req.Amount)
	if err != nil {
		log.Errorf("Failed emit: %v", err)
		resp.Error = &proto.Error{
			Code:    409,
			Message: "Can`t emit currency",
		}
	}

	resp.Balance, err = h.pgCl.GetUserBalances(req.UserId, req.CurrencyName)
	if err != nil {
		log.Errorf("Can`t get balance: %v", err)
	}

	respBody, err := marshalEmitBalanceByUserIdResponse(resp)
	if err != nil {
		log.Errorf("Failed serialize response: %v", err)
		return
	}

	h.sResponse.SendMessage("e.balances.forward", "r.balance.EmitUserBalanceResponse", *respBody)
	log.Infof("Send response for EmitBalanceByUserIdRequest, UserId [%s], ResponseId [%s]", resp.UserId, resp.Id)
}
