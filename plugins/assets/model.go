package PluginKits

import (
	"net/http"
)

type Config struct {
	CoreHandler CoreHandlerInterface
	RouterAPI   http.Handler
}
