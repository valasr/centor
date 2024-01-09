package grpc_server

import (
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type agentInfo struct {
	id       string  // id of the agent
	addr     string  // address of this node
	dc       string  // datacenter of this node
	isServer bool    // is this node a server or not
	isLeader bool    // is this node leader or not
	weight   int     // weight of this node in the cluster
	parent   *parent // parent of this node in the cluster or in primary cluster
	childs   map[string]*child
}
type agent struct {
	agentInfo
	// isPrimary bool   // is this node primary server or not
	isSubCluster bool      // is this node
	isReady      *brodBool // is this node ready or not
	isConneted   *brodBool // is this node connected to parent or not
	isStoping    *brodBool // is this node stoped or not

	listener   net.Listener
	grpcServer *grpc.Server
}

type clientStream struct {
	conn  *grpc.ClientConn      // connection to the server
	proto proto.DiscoveryClient // discovery protocol
	err   chan error            // channel for error
	close chan bool             // channel for closed connection
}

type parent struct {
	agentInfo    // parent agent information
	clientStream // stream of parent server
}

type child struct {
	agentInfo           // child agent information
	clientStream        // stream of the child server
	status       string // status of child in the cluster
}

// ===========================================
