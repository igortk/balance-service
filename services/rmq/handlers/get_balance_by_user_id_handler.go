package handlers

import (
	"balance-service/config"
	proto "balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type GetBalanceByUserIdHandler struct {
	pgCl      *pg.Client
	sResponse *senders.Sender
}

func NewGetBalanceByUserIdHandler(pgCl *pg.Client, s *senders.Sender) *GetBalanceByUserIdHandler {
	return &GetBalanceByUserIdHandler{
		pgCl:      pgCl,
		sResponse: s,
	}
}

func (h *GetBalanceByUserIdHandler) HandleMessage(body []byte) {
	req, err := unmarshalRequest[*proto.GetBalanceByUserIdRequest](body)
	if err != nil {
		log.Errorf("Failed deserialize request: %v", err)
		return
	}

	resp := h.processing(req)

	if err := h.send(resp); err != nil {
		log.Errorf("Can't send response: %v", err)
	}
}

func (h *GetBalanceByUserIdHandler) processing(req *proto.GetBalanceByUserIdRequest) *proto.GetBalanceByUserIdResponse {
	log.Infof("Start processing request by Id: %s", req.Id)

	resp := &proto.GetBalanceByUserIdResponse{Id: req.Id, UserId: req.UserId}

	reqErr := h.validation(req)
	if reqErr != nil {
		resp.Error = reqErr
		return resp
	}

	balances, err := h.pgCl.GetUserBalances(req.UserId)
	resp.UserBalance = balances

	if err != nil {
		log.Errorf("Can`t get user balances: %v", err)
		resp.Error = &proto.Error{
			Code:    409,
			Message: "Problem user balances",
		}
	}

	log.Infof("Finish processing request by Id: %s", req.Id)
	return resp
}

func (h *GetBalanceByUserIdHandler) validation(req *proto.GetBalanceByUserIdRequest) *proto.Error {
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
			Message: "Invalid user id:",
		}
	}

	return nil
}

func (h *GetBalanceByUserIdHandler) send(resp *proto.GetBalanceByUserIdResponse) error {
	respBody, err := marshalResponse(resp)
	if err != nil {
		return fmt.Errorf("failed serialize response GetBalanceByUserIdResponse: %v", err)
	}

	err = h.sResponse.SendMessage(config.RabbitBalanceExchange, config.GetBalanceByUserIdResponseRoutingKey, *respBody)
	if err != nil {
		return fmt.Errorf("failed send response GetBalanceByUserIdResponse: %v", err)
	}

	log.Infof("Send response for GetBalanceByUserIdRequest, UserId [%s], ResponseId [%s]", resp.UserId, resp.Id)
	return nil
}
