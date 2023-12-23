package grpc_server

import (
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type agent struct {
	id        string // id of the agent
	addr      string // address of this node
	dc        string // datacenter of this node
	isServer  bool   // is this node a server or not
	isLeader  bool   // is this node leader or not
	isPrimary bool   // is this node primary server or not
	isReady   bool   // is this node ready or not
	weight    int    // weight of this node in the cluster

	parent *parent           // parent of this node in the cluster or in primary cluster
	childs map[string]*child // childs of this node in the cluster
}

type stream struct {
	conn  *grpc.ClientConn      // connection to the server
	proto proto.DiscoveryClient // discovery protocol
	err   chan error            // channel for error
	close chan bool             // channel for closed connection
}

type parent struct {
	agent  // parent agent information
	stream // stream of parent server
}

type child struct {
	agent         // child agent information
	stream        // stream of the child server
	status string // status of child in the cluster
}

// ===========================================
