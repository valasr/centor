package grpc_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func server(ctx context.Context) (proto.DiscoveryClient, func()) {

	app, _ := grpc_server.NewServer(grpc_server.Config{
		Name:       "ali",
		DataCenter: "dc1",
		Host:       "localhost",
		Port:       3000,
		IsServer:   true,
		IsLeader:   true,
	})

	go func() {
		if err := app.Serve(nil); err != nil {
			log.Fatal(err)
		}
	}()

	app.GetCoreHandler().WaitForReady(ctx)

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	conn, err := grpc.DialContext(queryCtx, "localhost:3000",

		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	client := proto.NewDiscoveryClient(conn)

	closer := func() {

		fmt.Println("stopping.....")
		cancel()
		// if err := conn.Close(); err != nil {
		// 	log.Fatal(err)
		// }
		app.Stop()
	}

	return client, closer
}

func TestGetInfo(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *proto.InfoResponse
		err error
	}

	tests := map[string]struct {
		in       *proto.EmptyRequest
		expected expectation
	}{
		"Must_Success": {
			in: &proto.EmptyRequest{},
			expected: expectation{
				out: &proto.InfoResponse{
					Id: "ali",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetInfo(ctx, tt.in)
			if err != nil {
				t.Error("getInfo error :", err)

			} else {

				if tt.expected.out.Id != out.Id {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, out)
				}
			}

		})
	}
}
