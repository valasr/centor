package timeSyncer_plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrtdeh/centor/pkg/event"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
)

// ExamplePlugin1 is an example plugin implementing the Plugin interface
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

func (p *PluginProvider) Init() error {
	h = p.Handler

	r, ok := p.Router.(*gin.Engine)
	if !ok {
		return fmt.Errorf("router is not a gin router")
	}
	r.POST("/update-time", updateTimeRequest)

	p.Router = r

	err := event.Bus.Subscribe("sync-time", syncTime)
	if err != nil {
		return err
	}

	return nil
}

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)
}

type UpdateTimeRequest struct {
	NodeId     string `json:"node_id"`
	MainNodeId string `json:"main_node_id"`
}

func updateTimeRequest(c *gin.Context) {
	var req UpdateTimeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.FireEvent(context.Background(), req.NodeId, "sync-time", req.MainNodeId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "ok")
}

func syncTime(mainNode string) {
	fmt.Println("try to sync time with server... : ", mainNode)
}
