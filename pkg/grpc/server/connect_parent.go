package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

type connectConfig struct {
	ServersAddresses []string
}

func (a *agent) ConnectToParent(addrs []string) error {
	debug(a.id, "call ConnectToParent from %s", a.id)
	if len(addrs) == 0 {
		return nil
	}
	servers := addrs

	// master election for servers / best election for clients
	var si *ServerInfo
	var err error
	if a.isServer {
		// select leader only
		si, err = leaderElect(servers)
	} else {
		// select best server in server's pool
		si, err = bestElect(servers)
	}
	if err != nil {
		return err
	}

	// dial to selected server
	conn, err := grpc_Dial(DialConfig{
		Address: si.Addr,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	// create parent object
	a.parent = &parent{
		agentInfo: agentInfo{ // parent agent
			addr:     si.Addr,
			id:       si.Id,
			isLeader: si.IsLeader,
		},
		clientStream: clientStream{ // parent stream
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		},
	}
	if a.isLeader {
		if n, err := cluster.GetNode(a.id); err == nil {
			n.ParentId = si.Id
			cluster.UpdateNodes([]NodeInfo{*n})
		}
	}

	// create sync stream rpc to parent server
	go func() {
		err = grpc_Connect(context.Background(), a)
		if err != nil {
			a.parent.clientStream.err <- fmt.Errorf("error in sync : %s", err.Error())
		}
	}()

	// health check conenction for parent server
	go connHealthCheck(healthcheckOpt{
		Id:       fmt.Sprintf("%s-to-%s", a.id, a.parent.id),
		StopingC: a.isStoping,
		ClientS:  &a.parent.clientStream,
		Duration: time.Second * 5,
	})

	a.isConneted.Set(true)
	return <-a.parent.clientStream.err
}
