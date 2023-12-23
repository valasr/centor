package grpc_server

import (
	"context"
	"os/exec"

	"github.com/mrtdeh/centor/proto"
)

func (s *agent) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
	cmd := exec.Command("sh", "-c", req.Command)
	out, err := cmd.CombinedOutput()
	res := &proto.ExecResponse{
		Output: string(out),
	}
	if err != nil {

		return res, err
	}
	return res, nil
}
