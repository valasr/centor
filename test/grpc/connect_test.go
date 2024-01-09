package grpc_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/proto"
)

func listenServer(ctx context.Context) func() {
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
	app.GetCoreHandler().WaitForReady(ctx)

	closer := func() {
		fmt.Println("stopping.....")
		app.Stop()
	}

	return closer
}

func TestConnect(t *testing.T) {
	ctx := context.Background()

	closer := listenServer(ctx)
	defer closer()

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
				Name:       "reza",
				DataCenter: "dc1",
				IsServer:   true,
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

	var wg sync.WaitGroup
	for scenario, tt := range tests {
		wg.Add(1)
		t.Run(scenario, func(t *testing.T) {
			go func() {
				defer wg.Done()

				a, err := grpc_server.NewServer(*tt.in)
				if err != nil {
					t.Errorf("Err -> %s\n", err)
				}
				defer a.Stop()

				go func() {
					if err := a.Serve(nil); err != nil {
						// log.Fatal(err)
					}
				}()

				a.GetCoreHandler().WaitForConnect(ctx)

				res, err := a.Call(ctx, &proto.CallRequest{
					AgentId: a.GetCoreHandler().GetMyId(),
				})
				if err != nil {
					t.Errorf("Err -> \nGot: %q\n", err)
				} else {
					if len(res.Tags) != 2 {
						t.Errorf("tags must ne 2 but got %d\n", len(res.Tags))
					}

					fmt.Println("DEBUG : end")
					// time.Sleep(time.Second)
					// return
				}
			}()

			wg.Wait()
			fmt.Println("DEBUG : end 222")

		})
	}
	fmt.Println("DEBUG : end 333")
	// time.Sleep(time.Second * 3)
}
