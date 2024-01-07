package grpc_server

import (
	"fmt"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

func (a *agent) Serve(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	if lis == nil {
		var err error
		lis, err = net.Listen("tcp", a.addr)
		if err != nil {
			return fmt.Errorf("error creating the server %v", err)
		}
	}
	proto.RegisterDiscoveryServer(grpcServer, a)
	fmt.Println("listen an ", a.addr)

	a.listener = lis
	a.grpcServer = grpcServer

	a.isReady.Set(true)
	return grpcServer.Serve(lis)
}
