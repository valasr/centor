package grpc_server

import (
	"context"
	"errors"
	"sync"

	"github.com/mrtdeh/centor/pkg/kive"
	"github.com/mrtdeh/centor/proto"
)

type KVHandler struct{}

func (h *KVHandler) Sync(pr kive.PublishRequest) {
	k := KVPool{
		Key:    pr.Record.Key,
		Value:  pr.Record.Value,
		Action: pr.Action,
	}
	go kvm.SendKVtoAll(app, &k)
}

type KVPool struct {
	Key       string
	Value     string
	TargetId  string
	LastError string
	Action    string
	done      bool
}

type KVPoolManager struct {
	pools map[string]*KVPool
	l     sync.RWMutex
}

var kvm = &KVPoolManager{}

func GetKVManager() *KVPoolManager {
	return kvm
}

func (m *KVPoolManager) SendKVtoAll(a *agent, kvp *KVPool) error {
	m.l.Lock()
	defer m.l.Unlock()

	if a.parent != nil {
		tid := a.parent.id
		err := sendKVtoParent(a, kvp)
		if err != nil {
			kvp.TargetId = tid
			m.pools[tid] = kvp
			return nil
		}
	}
	if a.childs != nil {
		for _, c := range a.childs {
			tid := c.id
			err := sendKVtoChild(a, c, kvp)
			if err != nil {
				kvp.TargetId = tid
				m.pools[tid] = kvp
				return err
			}
		}
	}

	return nil
}

func sendKVtoChild(a *agent, c *child, kvp *KVPool) error {

	if c.status != StatusConnected {
		return errors.New("status is not connected")
	}
	res, err := c.proto.KVU(context.Background(), &proto.KVURequest{
		Key:    kvp.Key,
		Value:  kvp.Value,
		Action: kvp.Action,
		From:   a.id,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}

func sendKVtoParent(a *agent, kvp *KVPool) error {
	res, err := a.parent.proto.KVU(context.Background(), &proto.KVURequest{
		Key:    kvp.Key,
		Value:  kvp.Value,
		Action: kvp.Action,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

// func (m *KVPoolManager) Del(pid string) {
// 	m.l.Lock()
// 	defer m.l.Unlock()
// 	delete(m.pools, pid)
// }

func (a *agent) KVU(ctx context.Context, req *proto.KVURequest) (*proto.KVUResponse, error) {
	if req.Action == "add" {
		err := kive.Put(req.Key, req.Value)
		if err != nil {
			return &proto.KVUResponse{
				Error: err.Error(),
			}, nil
		}
	} else if req.Action == "delete" {
		err := kive.Del(req.Key)
		if err != nil {
			return &proto.KVUResponse{
				Error: err.Error(),
			}, nil
		}
	}

	return &proto.KVUResponse{}, nil
}
