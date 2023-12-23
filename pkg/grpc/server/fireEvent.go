package grpc_server

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mrtdeh/centor/pkg/event"
	"github.com/mrtdeh/centor/proto"
	goproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *agent) FireEvent(ctx context.Context, req *proto.EventRequest) (*proto.Close, error) {
	params := []interface{}{}
	for _, p := range req.Params {
		val, err := ConvertAnyToInterface(p)
		if err != nil {
			return nil, err
		}
		params = append(params, val)

	}

	if len(params) == 0 {
		event.Bus.Publish(req.Name)
	} else {
		event.Bus.Publish(req.Name, params...)
	}

	return &proto.Close{}, nil
}

func ConvertAnyToInterface(anyValue *any.Any) (interface{}, error) {
	var value interface{}
	bytesValue := &wrappers.BytesValue{}
	err := anypb.UnmarshalTo(anyValue, bytesValue, goproto.UnmarshalOptions{})
	if err != nil {
		return value, err
	}
	uErr := json.Unmarshal(bytesValue.Value, &value)
	if err != nil {
		return value, uErr
	}
	return value, nil
}
