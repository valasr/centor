package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mrtdeh/centor/proto"
)

const (
	MaxApplyTries = 1024
)

func (a *agent) syncAgentChange(ca *agent, action int32) error {
	ni := NodeInfo{
		Id:         ca.id,
		Address:    ca.addr,
		IsServer:   ca.isServer,
		IsLeader:   ca.isLeader,
		IsPrimary:  ca.isPrimary,
		DataCenter: ca.dc,
		ParentId:   ca.parent.id,
	}
	// if is leader then apply changes
	if a.isLeader {
		err := a.applyChange(ni, action)
		if err != nil {
			return fmt.Errorf("error in applyChange : %s", err.Error())
		}
	}
	// and if also is not primary then send change to primary
	if !a.isPrimary {
		var try int
		for {
			if try >= MaxApplyTries {
				return fmt.Errorf("max tries reached %d", MaxApplyTries)
			}
			try++

			if a.parent != nil {
				data, err := json.Marshal(ni)
				if err != nil {
					return err
				}
				_, err = a.parent.proto.Change(context.Background(), &proto.ChangeRequest{
					Change: &proto.ChangeRequest_NodesChange{
						NodesChange: &proto.NodesChange{
							Id:     ni.Id,
							Action: action,
							Data:   string(data),
						},
					},
				})
				if err != nil {
					log.Printf("error in syncChangeToLeader : %s", err.Error())
					time.Sleep(time.Second)
					continue
				}
				// break out of for
				break
			}
			// fmt.Println("retry to sync")
			time.Sleep(time.Millisecond * 100)
		}
	}

	return nil
}
