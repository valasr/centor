package grpc_server

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Call(ctx context.Context, req *proto.CallRequest) (*proto.CallResponse, error) {
	var tags []string
	tags = append(tags, a.id)

	if a.parent != nil && a.parent.id != req.AgentId {

		res, err := a.parent.proto.Call(context.Background(), &proto.CallRequest{
			AgentId: a.id,
		})
		if err != nil {
			return &proto.CallResponse{}, fmt.Errorf("failed to call parent %s: %v", a.parent.id, err)
		}
		tags = append(tags, res.Tags...)
	}

	if a.childs != nil {
		for _, c := range a.childs {
			if c.id != req.AgentId {
				res, err := c.proto.Call(context.Background(), &proto.CallRequest{
					AgentId: a.id,
				})
				if err != nil {
					return &proto.CallResponse{}, fmt.Errorf("failed to call child %s: %v", c.id, err)
				}
				tags = append(tags, res.Tags...)
			}
		}
	}

	return &proto.CallResponse{Tags: tags}, nil
}
