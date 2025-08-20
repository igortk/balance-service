package handlers

import (
	"fmt"
	gitProto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type MessageHandler interface {
	HandleMessage([]byte)
}

func unmarshalRequest[T proto.Message](body []byte) (T, error) {
	var zero T

	t := reflect.TypeOf(zero)
	if t == nil {
		return zero, fmt.Errorf("type is nil")
	}

	if t.Kind() != reflect.Ptr {
		return zero, fmt.Errorf("T must be a pointer to a proto.Message")
	}

	val := reflect.New(t.Elem()).Interface()

	msg, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("failed to cast created value to T")
	}

	if err := proto.Unmarshal(body, msg); err != nil {
		return zero, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
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
	return err != nil
}
