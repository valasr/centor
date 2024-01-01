package grpc_test

import (
	"context"
	"log"
	"net"
	"testing"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context) (proto.DiscoveryClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	app, _ := grpc_server.NewServer(grpc_server.Config{
		Name:       "ali",
		DataCenter: "dc1",
		Host:       "localhost",
		Port:       3000,
		IsServer:   true,
		IsLeader:   true,
	})

	go func() {
		if err := app.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		app.Stop()
	}

	client := proto.NewDiscoveryClient(conn)

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
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Id != out.Id {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}

		})
	}
}
