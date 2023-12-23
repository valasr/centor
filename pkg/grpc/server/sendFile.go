package grpc_server

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) SendFile(stream proto.Discovery_SendFileServer) error {
	var file *os.File
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			file.Close()
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading message: %w", err)
		}

		if info := msg.GetInfo(); info != nil {
			// in first, create package file to download from stream
			file, err = createFile(info)
			if err != nil {
				return fmt.Errorf("Error creating file: %w", err)
			}
		} else if chunk := msg.GetChunkData(); chunk != nil {
			// write chunk's package to created file
			file.Write(chunk)
		} else if end := msg.GetEnd(); end {
			// close created file
			fmt.Println("recieved new file : ", file.Name())
			file.Close()
			break
		}

	}
	return nil
}

func createFile(info *proto.FileInfo) (*os.File, error) {
	var err error
	if err := os.MkdirAll("/tmp/centor-recieved/", 0777); err != nil {
		return nil, err
	}
	filepath := path.Join("/tmp/centor-recieved/", path.Base(info.Name))
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
