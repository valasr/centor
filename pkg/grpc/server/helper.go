package grpc_server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	Aany "github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	goproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	Debug      bool       = true
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// ========================== HEALTH CHECK =========================
type healthcheckOpt struct {
	Id       string
	StopingC *brodBool
	ClientS  *clientStream
	Duration time.Duration
}

func connHealthCheck(opt healthcheckOpt) {
	msg := fmt.Sprintf("health check %s", opt.Id)

	debug("", "start %s", msg)
	defer func() {
		debug("", "closed %s", msg)
	}()

	timer := time.NewTimer(opt.Duration)

	for {
		select {
		case <-opt.StopingC.GetC():
			return
		case <-timer.C:
			if err := connIsFailed(opt.ClientS.conn); err != nil {
				opt.ClientS.err <- err
				return
			}
			_, err := opt.ClientS.proto.Ping(context.Background(), &proto.PingRequest{})
			if err != nil {
				opt.ClientS.err <- fmt.Errorf("error in ping request: %v", err)
				return
			}

			timer.Reset(opt.Duration)
		}

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

func getRandHash(length int) string {
	b := rand.Intn(10000)
	b64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", b)))

	if len(b64) < length {
		return b64
	}

	return b64[:length]
}

func debug(id, format string, a ...any) {
	if Debug {
		if id == "" {
			id = getRandHash(5)
		}
		id = "ID=" + id

		// fmt.Printf("[DEBUG]"+id+" :: "+format+"\n", a...)
		fmt.Printf("[DEBUG]::"+format+"::["+id+"]\n", a...)
	}
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
		fmt.Printf("Conenct Back from parent - ID=%s ME=%s\n", pid, a.id)
	}

}

func ConvertInterfaceToAny(v interface{}) (*Aany.Any, error) {
	anyValue := &Aany.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, goproto.MarshalOptions{})
	return anyValue, err
}
