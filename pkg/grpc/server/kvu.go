package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mrtdeh/centor/pkg/kive"
	"github.com/mrtdeh/centor/proto"
)

type KVUHandler struct {
	agent *agent
}

func (a *agent) GetKVUHandler() *KVUHandler {
	return &KVUHandler{
		agent: a,
	}
}
func (h *KVUHandler) Sync(pr *kive.PublishRequest, from string) {
	if pr != nil {
		go h.agent.sync(*pr, from)
	}
}

func (a *agent) sync(pr kive.PublishRequest, from string) {
	if from == "" {
		from = a.id
	}

	k := KVPool{
		Key:       pr.Record.Key,
		Value:     pr.Record.Value,
		Action:    pr.Action,
		Timestamp: pr.PublishDate,
		Namespace: pr.Namespace,
		From:      from,
	}
	go kvm.SendKVtoAll(a, &k)
}

type KVPool struct {
	Namespace string
	Key       string
	Value     string
	TargetId  string
	LastError string
	Action    string
	From      string
	Timestamp time.Time
	done      bool
}

type KVPoolManager struct {
	pools map[string]*KVPool
	l     sync.RWMutex
}

func (m *KVPoolManager) addPool(kvp *KVPool) {
	kvp.Timestamp = time.Now()
	m.pools[kvp.TargetId] = kvp
}

func (m *KVPoolManager) SendKVtoAll(a *agent, kvp *KVPool) error {
	m.l.Lock()
	defer m.l.Unlock()

	if a.parent != nil {
		tid := a.parent.id
		if tid != kvp.From {
			err := sendKVtoParent(a, kvp)
			if err != nil {
				kvp.TargetId = tid
				m.addPool(kvp)
				return nil
			}
		}
	}
	if a.childs != nil {
		for _, c := range a.childs {
			tid := c.id
			if tid != kvp.From {
				fmt.Println("send kv to child : ", tid)
				err := sendKVtoChild(a, c, kvp)
				if err != nil {
					log.Println("error in send kv to child : ", err)
					// kvp.TargetId = tid
					// m.addPool(kvp)
					return err
				}
			}
		}
	}

	return nil
}

var kvm = &KVPoolManager{
	pools: make(map[string]*KVPool),
}

func GetKVManager() *KVPoolManager {
	return kvm
}

func sendKVtoChild(a *agent, c *child, kvp *KVPool) error {

	if c.status != StatusConnected {
		return errors.New("status is not connected")
	}
	res, err := c.proto.KVU(context.Background(), &proto.KVURequest{
		Key:       kvp.Key,
		Value:     kvp.Value,
		Action:    kvp.Action,
		Namespace: kvp.Namespace,
		Timestamp: kvp.Timestamp.Unix(),
		From:      a.id,
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
		Key:       kvp.Key,
		Value:     kvp.Value,
		Action:    kvp.Action,
		Namespace: kvp.Namespace,
		Timestamp: kvp.Timestamp.Unix(),
		From:      kvp.From,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (a *agent) KVU(ctx context.Context, req *proto.KVURequest) (*proto.KVUResponse, error) {
	fmt.Println("hit kv update from ", req.From)
	var kv *kive.PublishRequest
	var err error

	if req.Action == "add" {
		kv, err = kive.Put(a.dc, req.Namespace, req.Key, req.Value, req.Timestamp)
		if err != nil {
			return &proto.KVUResponse{
				Error: err.Error(),
			}, nil
		}
	} else if req.Action == "delete" {
		kv, err = kive.Del(req.Namespace, req.Key, req.Timestamp)
		if err != nil {
			return &proto.KVUResponse{
				Error: err.Error(),
			}, nil
		}
	}

	if kv != nil {
		go a.sync(*kv, req.From)
	}
	return &proto.KVUResponse{}, nil
}
