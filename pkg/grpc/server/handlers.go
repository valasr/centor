package grpc_server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type CoreHandlers struct {
	agent *agent
}

func GetAgentHandler() *CoreHandlers {
	return &CoreHandlers{
		agent: app,
	}
}

type FileHandler struct {
	Name      string
	Extension string
	Data      []byte
}

// wait for current agent is running completely
func (h *CoreHandlers) WaitForReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if h.agent != nil && h.agent.isReady {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (h *CoreHandlers) GetMyId() string {
	h.WaitForReady(context.Background())
	return h.agent.id
}

func (h *CoreHandlers) GetParentId() string {
	h.WaitForReady(context.Background())
	if h.agent.parent != nil {
		return h.agent.parent.id
	}
	return ""
}

func (h *CoreHandlers) CallAPI(ctx context.Context, nodeId, method, addr, body string) (*map[string]interface{}, error) {
	if n, err := cluster.GetNode(nodeId); err == nil {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)

		// send event to server
		res, err := client.CallAPI(ctx, &proto.APIRequest{
			Method: method,
			Body:   body,
			Addr:   addr,
		})
		if err != nil {
			return nil, err
		}
		var result = &map[string]interface{}{
			"body":  res.Body,
			"error": res.Error,
		}

		return result, nil
	}

	return nil, fmt.Errorf("node id %s not found", nodeId)
}

func (h *CoreHandlers) FireEvent(ctx context.Context, nodeId, event string, params ...any) error {
	protoParams := []*anypb.Any{}

	var fire = func(addr string) error {
		conn, err := grpc_Dial(addr)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)

		// send event to server
		_, err = client.FireEvent(ctx, &proto.EventRequest{
			Name:   event,
			Params: protoParams,
			From:   h.agent.id,
		})
		if err != nil {
			return err
		}
		return nil
	}

	//  convert params to protobuf anypb
	for _, p := range params {
		anyValue, err := ConvertInterfaceToAny(p)
		if err != nil {
			return err
		}
		protoParams = append(protoParams, anyValue)
	}

	// first check node id with parent id
	if h.agent.isLeader && h.agent.parent != nil && h.agent.parent.id == nodeId {
		return fire(app.parent.addr)
	}

	// check if node id is exist in nodes or not
	if n, err := cluster.GetNode(nodeId); err == nil {
		return fire(n.Address)
	}
	return fmt.Errorf("node id %s is not exist", nodeId)
}

func (h *CoreHandlers) Exec(ctx context.Context, nodeId, commnad string) (string, error) {

	// check if node_id is exist or not
	if n, err := cluster.GetNode(nodeId); err == nil {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return "", err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)

		// run command on the connected server
		res, err := client.Exec(context.Background(), &proto.ExecRequest{
			Command: commnad,
		})
		if err != nil {
			return "", err
		}

		return res.Output, nil
	}

	return "", nil
}

func (h *CoreHandlers) SendFile(ctx context.Context, nodeId, filename string, data []byte) error {

	reader := bytes.NewReader(data)
	filesize := reader.Size()
	buffer := make([]byte, 1024)

	// check if node_id is exist or not
	if n, err := cluster.GetNode(nodeId); err == nil {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)
		stream, err := client.SendFile(context.Background())
		if err != nil {
			return err
		}

		// send the file information
		err = stream.Send(&proto.SendFileRequest{
			Data: &proto.SendFileRequest_Info{
				Info: &proto.FileInfo{
					Name: filename,
					Size: uint32(filesize),
				},
			},
		})
		if err != nil {
			return err
		}

		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			// send the chunks of the file
			err = stream.Send(&proto.SendFileRequest{
				Data: &proto.SendFileRequest_ChunkData{
					ChunkData: buffer[:n],
				},
			})
			if err != nil {
				return err
			}
		}

		// send the end of the file
		err = stream.Send(&proto.SendFileRequest{
			Data: &proto.SendFileRequest_End{
				End: true,
			},
		})
		if err != nil {
			return err
		}

		// receive server response and error if any
		_, err = stream.CloseAndRecv()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error receiving response: %v", err)
		}

	} else {
		return fmt.Errorf("Node %s not found", nodeId)
	}

	return nil
}

// Todo: this function should be removed
func (h *CoreHandlers) Call(ctx context.Context) (string, error) {

	res, err := h.agent.Call(ctx, &proto.CallRequest{
		AgentId: h.agent.id,
	})
	if err != nil {
		return "", err
	}
	return strings.Join(res.Tags, " ,"), nil
}

// returns a map of all the nodes in the cluster
func (h *CoreHandlers) GetClusterNodes() map[string]NodeInfo {
	return cluster.nodes
}
