package grpc_test

import (
	"context"
	"log"
	"testing"
	"time"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/proto"
)

func listenServer(ctx context.Context) {
	// buffer := 101024 * 1024
	// lis := bufconn.Listen(buffer)

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
}

func TestConnect(t *testing.T) {
	ctx := context.Background()

	listenServer(ctx)

	type expectation struct {
		out *proto.InfoResponse
		err error
	}

	tests := map[string]struct {
		in       *grpc_server.Config
		expected expectation
	}{
		"Must_Success": {
			in: &grpc_server.Config{
				Name:       "client-1",
				DataCenter: "dc1",
				Servers:    []string{"localhost:3000"},
				Host:       "localhost",
				Port:       3001,
			},
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

			go func() {
				a, err := grpc_server.NewServer(*tt.in)
				if err != nil {
					t.Errorf("Err -> %s\n", err)
				}

				err = a.ConnectToParent(tt.in.Servers)
				if err != nil {
					t.Errorf("Err -> \nGot: %q\n", err)
				} else {

					res, err := a.Call(context.Background(), &proto.CallRequest{
						AgentId: a.GetCoreHandler().GetMyId(),
					})
					if err != nil {
						t.Errorf("Err -> \nGot: %q\n", err)
					} else {
						if len(res.Tags) != 2 {
							t.Errorf("tags must ne 2 but got %d\n", len(res.Tags))
						}
					}

				}
			}()

		})
	}

	time.Sleep(time.Second * 3)
}
