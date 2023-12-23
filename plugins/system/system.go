package system_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrtdeh/centor/pkg/event"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
)

type PluginProvider struct {
	PluginKits.PluginProps
}

func (p *PluginProvider) SetHandler(h PluginKits.CoreHandlerInterface) {
	p.Handler = h
}

func (p *PluginProvider) SetRouter(r http.Handler) {
	p.Router = r
}

var h PluginKits.CoreHandlerInterface

var sysinfo system

type system struct {
	info map[string]System
	l    *sync.RWMutex
}

func (p *PluginProvider) Init() error {
	h = p.Handler
	sysinfo = system{
		info: make(map[string]System),
		l:    &sync.RWMutex{},
	}

	r, ok := p.Router.(*gin.Engine)
	if !ok {
		return fmt.Errorf("router is not a gin router")
	}
	r.GET("/system-info", getSystemInfo)

	p.Router = r

	err := event.Bus.Subscribe("system-info", updateSystemsInfoEvent)
	if err != nil {
		return err
	}

	return nil
}

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)

	h.WaitForReady(context.Background())
	fmt.Println("debug system info")
	i := 10
	for {

		err := updateSystemsInfo([]System{
			{
				NodeId:     h.GetMyId(),
				DataCenter: "dc123",
				Hostname:   "myhost",
				RAM:        fmt.Sprintf("%d%%", i),
				CPU:        "30%",
				Disk:       "80%",
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 1)
		i++
	}
}

type System struct {
	RAM        string `json:"ram"`
	CPU        string `json:"cpu"`
	Disk       string `json:"disk"`
	Hostname   string `json:"host"`
	DataCenter string `json:"dc"`
	NodeId     string `json:"id"`
}

func updateSystemsInfo(systems []System) error {
	sysinfo.l.Lock()
	for _, s := range systems {
		sysinfo.info[s.NodeId] = s
	}
	sysinfo.l.Unlock()

	parentId := h.GetParentId()
	if parentId != "" {
		data, err := json.Marshal(infoToArray())
		if err != nil {
			return err
		}
		err = h.FireEvent(context.Background(), parentId, "system-info", h.GetMyId(), string(data))
		if err != nil {
			return err
		}
	}
	return nil
}

func updateSystemsInfoEvent(nodeId string, info string) {
	var systems []System
	err := json.Unmarshal([]byte(info), &systems)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = updateSystemsInfo(systems)
	if err != nil {
		fmt.Println("error in update event : ", err)
	}
}

func infoToArray() []System {
	sysinfo.l.Lock()
	var res = []System{}
	for _, system := range sysinfo.info {
		res = append(res, system)
	}
	sysinfo.l.Unlock()
	return res
}

func getSystemInfo(c *gin.Context) {
	c.JSON(200, infoToArray())
}
