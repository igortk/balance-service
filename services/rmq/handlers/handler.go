package handlers

import (
	"balance-service/dto/proto"
	"fmt"
	gitProto "github.com/golang/protobuf/proto"
)

type MessageHandler interface {
	HandleMessage([]byte)
}

func unmarshalEmitBalanceByUserIdRequest(body *[]byte) (*proto.EmitBalanceByUserIdRequest, error) {
	req := &proto.EmitBalanceByUserIdRequest{}

	err := gitProto.Unmarshal(*body, req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return req, nil
}

func marshalEmitBalanceByUserIdResponse(resp *proto.EmitBalanceByUserIdResponse) (*[]byte, error) {
	body, err := gitProto.Marshal(resp)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return &body, nil
}
