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

	a.listener = lis
	a.grpcServer = grpcServer
	a.isReady.Set(true)

	debug(a.id, "listen an %s", a.addr)
	return grpcServer.Serve(lis)
}
