package handlers

import (
	"balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var minAmount = 0.1

type EmitBalanceByUserIdHandler struct {
	pgCl      *pg.Client
	sResponse *senders.Sender
}

func NewEmitBalanceByUserIdHandler(pgCl *pg.Client, s *senders.Sender) *EmitBalanceByUserIdHandler {
	return &EmitBalanceByUserIdHandler{
		pgCl:      pgCl,
		sResponse: s,
	}
}

func (h *EmitBalanceByUserIdHandler) HandleMessage(body []byte) {
	req, err := unmarshalRequest[*proto.EmitBalanceByUserIdRequest](body)
	if err != nil {
		log.Errorf("Failed deserialize EmitBalanceByUserIdRequest: %v", err)
		return
	}

	resp := h.processing(req)

	if err := h.send(resp); err != nil {
		log.Errorf("Can't send response: %v", err)
	}
}

func (h *EmitBalanceByUserIdHandler) processing(req *proto.EmitBalanceByUserIdRequest) *proto.EmitBalanceByUserIdResponse {
	log.Infof("Start processing request by Id: %s", req.Id)

	resp := &proto.EmitBalanceByUserIdResponse{Id: req.Id, UserId: req.UserId}
	reqErr := h.validation(req)
	if reqErr != nil {
		resp.Error = reqErr
		return resp
	}

	err := h.pgCl.EmitCurrency(req.UserId, req.CurrencyName, req.Amount)
	if err != nil {
		log.Errorf("Failed emit: %v", err)
		resp.Error = &proto.Error{
			Code:    409,
			Message: "Can`t emit currency",
		}
		return resp
	}

	resp.Balance, err = h.pgCl.GetUserBalance(req.UserId, req.CurrencyName)
	if err != nil {
		log.Errorf("Can`t get balance: %v", err)
		resp.Error = &proto.Error{
			Code:    409,
			Message: "Problem user balance",
		}
		return resp
	}

	log.Infof("Finish processing request by Id: %s", req.Id)
	return resp
}

func (h *EmitBalanceByUserIdHandler) validation(req *proto.EmitBalanceByUserIdRequest) *proto.Error {
	if isValidUUID(req.Id) {
		log.Errorf("Invalid request id: %v", req.Id)
		return &proto.Error{
			Code:    409,
			Message: "Invalid request id",
		}
	}

	if isValidUUID(req.UserId) {
		log.Errorf("Invalid user id: %v", req.UserId)
		return &proto.Error{
			Code:    409,
			Message: "Invalid user id",
		}
	}

	if req.Amount < minAmount {
		log.Errorf("Invalid user id: %v", req.UserId)
		return &proto.Error{
			Code:    409,
			Message: fmt.Sprintf("Amount less then %f", minAmount),
		}
	}

	return nil
}

func (h *EmitBalanceByUserIdHandler) send(resp *proto.EmitBalanceByUserIdResponse) error {
	respBody, err := marshalResponse(resp)
	if err != nil {
		return fmt.Errorf("failed serialize response EmitBalanceByUserIdResponse: %v", err)
	}

	err = h.sResponse.SendMessage("e.balances.forward", "r.balance-service.EmitUserBalanceResponse", *respBody)
	if err != nil {
		return fmt.Errorf("failed send response EmitBalanceByUserIdResponse: %v", err)
	}
	log.Infof("Send response for EmitBalanceByUserIdRequest, UserId [%s], ResponseId [%s]", resp.UserId, resp.Id)

	return nil
}
