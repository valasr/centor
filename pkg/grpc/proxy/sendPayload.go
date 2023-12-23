package grpc_proxy

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
)

// client read payload bytes from other client
func (sp *server) SendPayload(ctx context.Context, req *proto.RequestPayload) (*proto.ResponsePayload, error) {
	// log.Println("new payload : ", string(req.Conn))
	sp.msgIn <- req.Conn

	select {

	case res := <-sp.msgOut:
		return &proto.ResponsePayload{
			Body: res,
		}, nil
	case err := <-sp.err:
		log.Println("error in send payload : ", err.Error())
		return nil, err
	}

}

func (sp *server) tcpDialToService(port string) {

	for {
		fmt.Println("connecting to localhost:", port)

		dst, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", port))
		if err != nil {
			panic("Dial Error:" + err.Error())
		}
		defer dst.Close()

		go func() {
			select {
			case input := <-sp.msgIn:
				fmt.Println("debug input : ", string(input[:16]))
				dst.Write(input)
			}
			fmt.Println("debug end input")

		}()

		connData := make([]byte, 1024)
		_, err = dst.Read(connData)
		if err != nil {
			log.Println("error in read connection : ", err.Error())
		}
		sp.msgOut <- connData
		dst.Close()
	}

}

// func copyConn(src, dst, out bytes.Buffer) {
// 	// defer dst.Close()

// 	go func() {
// 		io.Copy(dst, src)
// 	}()

// 	io.Copy(out, dst)
// }
