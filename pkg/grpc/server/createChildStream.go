package grpc_server

import (
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) CreateChildStream(c *child, done chan bool) error {
	// dial to child listener
	conn, err := grpc_Dial(DialConfig{
		Address: c.addr,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	// create child object
	if cc, ok := a.childs[c.id]; ok {
		// store client connection and proto info
		cc.clientStream = clientStream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		}

		done <- true
		// run health check conenction for this child

		go connHealthCheck(healthcheckOpt{
			Id:       fmt.Sprintf("%s-to-%s", a.id, cc.id),
			StopingC: a.isStoping,
			ClientS:  &cc.clientStream,
			Duration: time.Second * 5,
		})

	} else {
		return fmt.Errorf("child you want to check not exist")
	}

	// status update for child
	c.status = StatusConnected

	// return back error message when child is disconnected or failed
	return <-c.clientStream.err
}
