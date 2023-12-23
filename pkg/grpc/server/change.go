package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mrtdeh/centor/proto"
)

var (
	cluster *ClusterInfo
)

type (
	NodeInfo struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
		Port       string `json:"port"`
		IsServer   bool   `json:"is_server"`
		IsLeader   bool   `json:"is_leader"`
		IsPrimary  bool   `json:"is_primary"`
		ParentId   string `json:"parent_id"`
		DataCenter string `json:"data_center"`
	}
	NodesInfoMap map[string]NodeInfo
	ClusterInfo  struct {
		nodes NodesInfoMap
		l     *sync.RWMutex
	}
)

func (n *NodesInfoMap) toArray() (narr []NodeInfo) {
	for _, v := range *n {
		narr = append(narr, v)
	}
	return narr
}

func (c *ClusterInfo) DeleteNode(nodeId string) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.nodes, nodeId)
}

func (c *ClusterInfo) UpdateNodes(nodes []NodeInfo) {
	c.l.Lock()
	defer c.l.Unlock()
	for _, node := range nodes {
		cluster.nodes[node.Id] = node
	}
}

func (c *ClusterInfo) GetNode(nodeId string) (*NodeInfo, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	if n, ok := cluster.nodes[nodeId]; ok {
		return &n, nil
	}
	return nil, fmt.Errorf("node id not found in cluster nodes map")
}

func init() {
	cluster = &ClusterInfo{
		nodes: make(map[string]NodeInfo),
		l:     &sync.RWMutex{},
	}
}

func (a *agent) Change(ctx context.Context, req *proto.ChangeRequest) (*proto.Close, error) {
	c := &proto.Close{}
	if !a.isLeader {
		return c, fmt.Errorf("you must send change request to primary not here")
	}

	if nch := req.GetNodesChange(); nch != nil {
		// fmt.Println("New change - change nodes")

		var ni NodeInfo
		err := json.Unmarshal([]byte(nch.Data), &ni)
		if err != nil {
			return c, err
		}

		err = a.syncAgentChange(&agent{
			id:        ni.Id,
			addr:      ni.Address,
			isServer:  ni.IsServer,
			isLeader:  ni.IsLeader,
			isPrimary: ni.IsPrimary,
			dc:        ni.DataCenter,
			parent:    &parent{agent: agent{id: ni.ParentId}},
		}, nch.Action)
		if err != nil {
			return c, err
		}

	}

	return c, nil
}
