package grpc_server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *agent) CreateProxy(ctx context.Context, req *proto.CreateProxyRequest) (*proto.Close, error) {
	// search service id in services pool and get that target id

	// iterate connections to find target
	for _, v := range s.childs {
		if req.TargetId == v.Id {

			// send proxy request to target
			v.conn.Send(&proto.LeaderResponse{
				Data: &proto.LeaderResponse_ProxyRequest{
					ProxyRequest: &proto.ProxyRequest{
						ProxyPort: req.TargetServicePort,
					},
				},
			})

			break
		}
	}
	return &proto.Close{}, nil
}

func (a *agent) CreateServiceProxy() {

	// check for ready
	a.waitForReady()

	// create proxy request
	if err := a.createProxy("service id"); err != nil {
		log.Fatalf("error in proxy : %s", err.Error())
	}

	// connect to proxy server (gRPC)
	cc, err := grpc.Dial("localhost:11001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error in dial : %s", err.Error())
	}
	ccc := proto.NewProxyManagerClient(cc)

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic("connection error:" + err.Error())
	}

	for {
		lc, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}
		// copyConn(conn)
		connData := make([]byte, 1024)
		_, err = lc.Read(connData)
		if err != nil {
			if err == io.EOF {
				log.Println("read EOF")
				continue
			}
			log.Fatal("error in read connection222 : ", err.Error())
		}
		res, err := ccc.SendPayload(context.Background(), &proto.RequestPayload{
			Conn: connData,
		})
		if err != nil {
			log.Fatal("error in SendPayload : ", err.Error())
		}

		lc.Write(res.Body)
	}

}

func (c *agent) createProxy(serviceId string) error {
	_, err := c.parent.conn.CreateProxy(context.Background(), &proto.CreateProxyRequest{
		// ServiceId: serviceId,
		// TargetId:          target,
		// TargetServicePort: port,
	})
	if err != nil {
		return fmt.Errorf("failed to proxy : %s\n", err.Error())
	}

	return nil
}
