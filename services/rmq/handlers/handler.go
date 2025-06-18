package handlers

import (
	"fmt"
	gitProto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

type MessageHandler interface {
	HandleMessage([]byte)
}

func unmarshalRequest[T gitProto.Message](body *[]byte) (T, error) {
	var req T

	err := gitProto.Unmarshal(*body, req)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return req, nil
}

func marshalResponse(resp gitProto.Message) (*[]byte, error) {
	body, err := gitProto.Marshal(resp)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return &body, nil
}

func isValidUUID(field string) bool {
	_, err := uuid.Parse(field)
	return err == nil
}
