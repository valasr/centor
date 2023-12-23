package grpc_server

import (
	"fmt"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

func (a *agent) Listen() error {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, a)
	fmt.Println("listen an ", a.addr)
	a.isReady = true
	return grpcServer.Serve(listener)
}
