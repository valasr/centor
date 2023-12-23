package grpc_server

import (
	"context"
	"fmt"
	"math"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerInfo struct {
	Id       string
	Addr     string
	IsLeader bool
}

func bestElect(addrs []string) (*ServerInfo, error) {

	index := -1
	weight := math.MaxInt32
	var si *ServerInfo
	for i, a := range addrs {

		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("error in dial : %s", err.Error())
		}

		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			fmt.Println("failed to get info from  :", a)
			continue
		}

		if res.Weight < int32(weight) {
			weight = int(res.Weight)
			index = i
			si = &ServerInfo{res.Id, a, res.IsLeader}
			conn.Close()
		}

	}

	if index == -1 {
		return nil, fmt.Errorf("server's are not available")
	}

	return si, nil

}

func leaderElect(addrs []string) (*ServerInfo, error) {
	var si *ServerInfo
	for _, a := range addrs {
		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("error in dial : %s", err.Error())
		}
		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			return nil, fmt.Errorf("error in getInfo : %s", err.Error())
		}

		si = &ServerInfo{res.Id, a, res.IsLeader}

		if res.IsLeader {
			conn.Close()
			return si, nil
		}

	}
	return nil, fmt.Errorf("leader not found")
}
