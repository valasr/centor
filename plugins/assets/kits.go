package PluginKits

import (
	"context"
	"fmt"
	"net/http"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

type PluginProps struct {
	Name    string
	Handler CoreHandlerInterface
	Router  http.Handler
}

type CoreHandlerInterface interface {
	Call(context.Context) (string, error)
	GetClusterNodes() map[string]grpc_server.NodeInfo
	WaitForReady(context.Context) error
	SendFile(context.Context, string, string, []byte) error
	Exec(context.Context, string, string) (string, error)
	FireEvent(context.Context, string, string, ...any) error
	CallAPI(context.Context, string, string, string, string) (*map[string]interface{}, error)
	GetParentId() string
	GetMyId() string
}

// Plugin interface
type Plugin interface {
	Init() error
	Run()
	SetHandler(handler CoreHandlerInterface)
	SetRouter(router http.Handler)
}

// CoreService structure
type PluginManagerService struct {
	Plugins []Plugin
}

// AddPlugin method to add a plugin to the CoreService
func (c *PluginManagerService) AddPlugin(p Plugin) {
	c.Plugins = append(c.Plugins, p)
}

// Start method to start the CoreService and its plugins
func (c *PluginManagerService) Start(h CoreHandlerInterface, r http.Handler) {
	fmt.Println("PluginManagerService is starting...")
	for _, plugin := range c.Plugins {
		plugin.SetHandler(h)
		plugin.SetRouter(r)
		if err := plugin.Init(); err != nil {
			fmt.Println("error in initializing plugin : ", err)
			continue
		}
		go plugin.Run()
	}
}
