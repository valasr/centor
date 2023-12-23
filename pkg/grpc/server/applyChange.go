package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

const (
	ChangeActionAdd = iota
	ChangeActionRemove
)

func (a *agent) applyChange(ni NodeInfo, action int32) error {
	id := ni.Id
	if id == "" {
		return fmt.Errorf("id is empty, must be exist")
	}
	switch action {
	case ChangeActionAdd:
		cluster.UpdateNodes([]NodeInfo{ni})
	case ChangeActionRemove:
		cluster.DeleteNode(id)
	}

	data, err := json.Marshal(cluster.nodes)
	if err != nil {
		return err
	}

	for _, child := range a.childs {
		// ignore disconnected child
		if child.status != StatusConnected {
			fmt.Println("apply change to child id : ", child.id, " status : ", child.status)
			continue
		}
		// ignore leader child
		if child.isLeader {
			continue
		}
		// notice to child
		_, err := child.proto.Notice(context.Background(), &proto.NoticeRequest{
			Notice: &proto.NoticeRequest_NodesChange{
				NodesChange: &proto.NodesChange{
					Id:   a.id,
					Data: string(data),
				},
			},
		})
		if err != nil {
			return err
		}

	}

	return nil
}
