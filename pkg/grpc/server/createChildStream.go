package grpc_server

import (
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) CreatechildStream(c *child, done chan bool) error {
	// dial to child listener
	conn, err := grpc_Dial(c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// create child object
	if cc, ok := a.childs[c.id]; ok {
		// store client connection and proto info
		cc.stream = stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		}

		done <- true
		// run health check conenction for this child
		go connHealthCheck(&cc.stream, time.Second*2)
	} else {
		return fmt.Errorf("child you want to check not exist")
	}

	// status update for child
	c.status = StatusConnected

	// return back error message when child is disconnected or failed
	return <-c.stream.err
}
