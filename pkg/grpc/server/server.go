package grpc_server

import (
	"fmt"
	"time"
)

type Config struct {
	Name       string   // Name of the server(id)
	DataCenter string   // Name of the server(id)
	Host       string   // Host of the server
	AltHost    string   // alternative host of the server (optional)
	Port       uint     // Port of the server
	Servers    []string // servers addresses for replication
	Primaries  []string // primaries addresses
	IsServer   bool     // is this node a server or not
	IsLeader   bool     // is this node leader or not
}

func (a *agent) Stop() error {
	fmt.Println("Stopping agent...")
	return nil
}

func NewServer(cnf Config) (*agent, error) {
	if cnf.Host == "" {
		cnf.Host = "127.0.0.1"
	}

	a := &agent{}

	// resolve alternative host from config
	var host string = cnf.Host
	if cnf.AltHost != "" {
		host = cnf.AltHost
	}
	// create default agent instance
	*a = agent{
		id:       cnf.Name,
		dc:       cnf.DataCenter,
		addr:     fmt.Sprintf("%s:%d", host, cnf.Port),
		childs:   make(map[string]*child),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
	}

	if !cnf.IsLeader || len(cnf.Primaries) > 0 {
		// if this node is a leader and no primaries are specified, this node becomes primary
		a.isSubCluster = true
	}

	var servers []string
	var connectToPrimary bool

	if !cnf.IsLeader && len(cnf.Servers) > 0 {
		servers = cnf.Servers
	} else {
		// if is a leader or there are no servers in the cluster
		// add current node info to nodes info map
		cluster.UpdateNodes([]NodeInfo{
			{
				Id:       a.id,
				Address:  a.addr,
				IsServer: a.isServer,
				IsLeader: a.isLeader,
				// IsPrimary:  a.isPrimary,
				DataCenter: a.dc,
			},
		})

		if len(cnf.Primaries) > 0 {
			connectToPrimary = true
			servers = cnf.Primaries
		}
	}

	go func() {
		for {
			// try connect to parent server
			err := a.ConnectToParent(connectConfig{
				ConnectToPrimary: connectToPrimary,
				ServersAddresses: servers,
			})
			if err != nil {
				fmt.Println(err.Error())
			}
			// retry delay time 1 second
			time.Sleep(time.Second * 1)
		}
	}()

	fmt.Println("DataCenter : ", cnf.DataCenter)
	return a, nil
}
