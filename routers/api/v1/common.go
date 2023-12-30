package api_v1

import grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"

var h *grpc_server.CoreHandlers

func Init(serverHandler *grpc_server.CoreHandlers) error {
	h = serverHandler
	return nil
}
