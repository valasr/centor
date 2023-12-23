package apiCaller_plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	r.POST("/call-api", callAPI)

	p.Router = r

	return nil
}

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)
}

type APIRequest struct {
	Address string `json:"address"`
	Method  string `json:"method"`
	Body    string `json:"body"`
	NodeId  string `json:"node_id"`
}

func callAPI(c *gin.Context) {
	var req APIRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.CallAPI(context.Background(), req.NodeId, req.Method, req.Address, req.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}
