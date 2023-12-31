package main

import (
	"fmt"
	"log"
	"strings"

	api_server "github.com/mrtdeh/centor/pkg/api"
	"github.com/mrtdeh/centor/pkg/cli"
	"github.com/mrtdeh/centor/pkg/config"
	"github.com/mrtdeh/centor/pkg/envoy"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/pkg/kive"
	pluginManager "github.com/mrtdeh/centor/plugins"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
	"github.com/mrtdeh/centor/routers"
	api_v1 "github.com/mrtdeh/centor/routers/api/v1"
)

func main() {

	// print centor in cli
	cli.PrintLogo()

	// load configurations
	cnf := config.LoadConfiguration()

	var serversAddrs []string
	sd := cnf.ServersAddr
	if sd != "" {
		serversAddrs = strings.Split(strings.TrimSpace(sd), ",")
	}

	var primariesAddrs []string
	pd := cnf.PrimaryServersAddr
	if pd != "" {
		primariesAddrs = strings.Split(strings.TrimSpace(pd), ",")
	}

	// create new gRPC server instance
	app, _ := grpc_server.NewServer(grpc_server.Config{
		Name:       cnf.Name,
		DataCenter: cnf.DataCenter,
		Host:       cnf.Host,
		AltHost:    cnf.AltHost,
		Port:       cnf.Port,
		IsServer:   cnf.IsServer,
		IsLeader:   cnf.IsLeader,
		Servers:    serversAddrs,
		Primaries:  primariesAddrs,
	})

	// initilize api server
	router := routers.InitRouter()
	api_v1.Init(app.GetCoreHandler())
	kive.Init(app.GetKVUHandler())

	// bootstrap plugins
	err := pluginManager.Bootstrap(pluginManager.Config{
		Config: PluginKits.Config{
			CoreHandler: app.GetCoreHandler(),
			RouterAPI:   router,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// start api server
	if config.WithAPI {
		httpServer := api_server.HttpServer{
			Host:   "0.0.0.0",
			Port:   9090,
			Router: router,
		}
		fmt.Printf("initil api server an address %s:%d\n", httpServer.Host, httpServer.Port)
		go func() {
			log.Fatal(httpServer.Serve())
		}()
	}
	// start envoy proxy server
	if cnf.Connect != "" {
		go func() {
			log.Fatal(envoy.NewEnvoy(
				envoy.EnvoyConfig{
					// note: give container name if running on docker
					EndpointAddress: cnf.AltHost,
					ListenerPort:    cnf.Port + 1,
					EndpointPort:    cnf.Port,
					TLSConfig: envoy.TLSConfig{
						Secure:         cnf.SSL_Enabled,
						CA:             cnf.SSL_ca,
						Cert:           cnf.SSL_cert,
						Key:            cnf.SSL_key,
						SessionTimeout: "6000s",
					},
				},
			))
		}()
		log.Println("start envoy server at :", cnf.Port+1)
	}

	// start gRPC server
	if err := app.Serve(nil); err != nil {
		log.Fatal(err)
	}

}
