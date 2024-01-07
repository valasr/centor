package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	goproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (a *agent) parentErr() <-chan error {
	return a.parent.clientStream.err
}

// =======================================

func (c *child) childErr() <-chan error {
	return c.clientStream.err
}

// ======================================

func (a *agent) Closechild(c *child) error {
	if _, ok := a.childs[c.id]; !ok {
		return fmt.Errorf("child %s is not exist", c.id)
	}
	delete(a.childs, c.id)
	return nil
}

// ========================== HEALTH CHECK =========================
func connHealthCheck(s *clientStream, d time.Duration) {
	for {
		if err := connIsFailed(s.conn); err != nil {
			s.err <- err
			return
		}
		_, err := s.proto.Ping(context.Background(), &proto.PingRequest{})
		if err != nil {
			s.err <- err
			return
		}
		time.Sleep(d)
	}
}

func connIsFailed(conn *grpc.ClientConn) error {
	status := conn.GetState()
	if status == connectivity.TransientFailure ||
		// status == connectivity.Idle ||
		status == connectivity.Shutdown {
		return fmt.Errorf("connection is failed with status %s", status)
	}
	return nil
}

// =============================================================
type DialConfig struct {
	Address     string
	DialContext func(context.Context, string) (net.Conn, error)
}

func grpc_Dial(cnf DialConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if cnf.DialContext != nil {
		cnf.Address = ""
		opts = append(opts, grpc.WithContextDialer(cnf.DialContext))
	}
	conn, err := grpc.Dial(cnf.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("error in dial : %s", err.Error())
	}
	return conn, nil
}

func grpc_Connect(ctx context.Context, a *agent) error {
	stream, err := a.parent.proto.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error in create connect stream : %s", err.Error())
	}

	// send connect message to parent server
	err = stream.Send(&proto.ConnectMessage{
		Id:         a.id,
		DataCenter: a.dc,
		Addr:       a.addr,
		IsServer:   a.isServer,
		IsLeader:   a.isLeader,
		ParentId:   a.parent.id,
	})
	if err != nil {
		return fmt.Errorf("error in send connect message : %s", err.Error())
	}

	var pid string
	for {
		// receive connect message from parent server
		res, err := stream.Recv()
		if err != nil {
			a.parent.clientStream.err <- err
			if pid != "" {
				fmt.Printf("Disconnect parent - ID=%s\n", pid)
			} else {
				return fmt.Errorf("error in receive connect message : %s", err.Error())
			}
			return err
		}
		if res != nil {
			pid = res.Id
		}
		fmt.Printf("Conenct Back from parent - ID=%s\n", pid)
	}

}

func ConvertInterfaceToAny(v interface{}) (*any.Any, error) {
	anyValue := &any.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, goproto.MarshalOptions{})
	return anyValue, err
}
