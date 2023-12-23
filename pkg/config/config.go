package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type config struct {
	Name               *string `hcl:"name"`
	DataCenter         *string `hcl:"dc"`
	Host               *string `hcl:"host" json:"host,omitempty"`
	Port               *uint   `hcl:"port"`
	IsServer           *bool   `hcl:"is_server"`
	IsLeader           *bool   `hcl:"is_leader"`
	ServersAddr        *string `hcl:"servers_address"`
	PrimaryServersAddr *string `hcl:"primaries_address"`
	Service            *struct {
		Id   string `hcl:"id"`
		Name string `hcl:"name"`
		Port uint   `hcl:"port"`
	} `hcl:"service,block"`
}

type Service struct {
	Id   string
	Name string
	Port uint
}

type Config struct {
	Name               string    // id of the agent
	DataCenter         string    // datacenter of the agent
	Host               string    // hostname of the agent
	AltHost            string    // hostname of the agent (alternative) (optional)
	Port               uint      // port of the agent
	IsServer           bool      // is this node a server or not
	IsLeader           bool      // is this node leader or not
	ServersAddr        string    // address of the servers in the cluster
	PrimaryServersAddr string    // address of the primary servers in the cluster
	Services           []Service // services in the cluster
	Connect            string

	SSL_Enabled bool
	SSL_ca      string
	SSL_cert    string
	SSL_key     string
}

var (
	Verbose bool // verbose mode
	WithAPI bool // with api endpoint
)

func LoadConfiguration() *Config {

	path := "/etc/centor.d/"

	// load configuration from file
	configs, err := loadConfigsFromDir(path)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// compile hcl configuration to struct
	cnf, err := compile(configs)
	if err != nil {
		log.Fatalf("Failed to compile configuration: %s", err)
	}

	// get environment variables
	if e := os.Getenv("PORT"); e != "" {
		p, _ := strconv.Atoi(e)
		cnf.Port = uint(p)
	}
	if e := os.Getenv("NAME"); e != "" {
		cnf.Name = e
	}
	if e := os.Getenv("DC"); e != "" {
		cnf.DataCenter = e
	}
	if e := os.Getenv("CONNECT"); e != "" {
		cnf.Connect = e
	}
	if e := os.Getenv("HOST"); e != "" {
		cnf.Host = e
	}
	if e := os.Getenv("JOIN"); e != "" {
		cnf.ServersAddr = e
	}
	if e := os.Getenv("SSL_ENABLED"); isTrue(e) {
		cnf.SSL_Enabled = true
	}
	if e := os.Getenv("SSL_CA"); e != "" {
		cnf.SSL_ca = e
	}
	if e := os.Getenv("SSL_CERT"); e != "" {
		cnf.SSL_cert = e
	}
	if e := os.Getenv("SSL_KEY"); e != "" {
		cnf.SSL_key = e
	}
	if e := os.Getenv("PRIMARIES"); e != "" {
		cnf.PrimaryServersAddr = e
	}
	if e := os.Getenv("SERVER"); isTrue(e) {
		cnf.IsServer = true
	}
	if e := os.Getenv("LEADER"); isTrue(e) {
		cnf.IsLeader = true
	}
	if e := os.Getenv("ALTERNATIVE_HOST"); e != "" {
		cnf.AltHost = e
	}

	// load config from cli arguments
	flag.BoolVar(&Verbose, "v", false, "")
	flag.BoolVar(&WithAPI, "api", false, "")

	flag.StringVar(&cnf.Name, "n", cnf.Name, "")
	flag.StringVar(&cnf.DataCenter, "dc", cnf.DataCenter, "")

	flag.BoolVar(&cnf.SSL_Enabled, "ssl_enabled", cnf.SSL_Enabled, "")
	flag.StringVar(&cnf.SSL_ca, "ssl_ca", cnf.SSL_ca, "")
	flag.StringVar(&cnf.SSL_cert, "ssl_cert", cnf.SSL_cert, "")
	flag.StringVar(&cnf.SSL_key, "ssl_key", cnf.SSL_key, "")

	flag.StringVar(&cnf.Host, "h", cnf.Host, "")
	flag.StringVar(&cnf.AltHost, "ah", cnf.AltHost, "")
	flag.UintVar(&cnf.Port, "p", cnf.Port, "")

	flag.StringVar(&cnf.Connect, "connect", cnf.Connect, "")
	flag.StringVar(&cnf.PrimaryServersAddr, "primaries-addr", cnf.PrimaryServersAddr, "")
	flag.StringVar(&cnf.ServersAddr, "join", cnf.ServersAddr, "")

	flag.BoolVar(&cnf.IsServer, "server", cnf.IsServer, "")
	flag.BoolVar(&cnf.IsLeader, "leader", cnf.IsLeader, "")
	flag.Parse()

	// print configuration an verbose mode
	if Verbose {
		cb, _ := json.MarshalIndent(cnf, "", " ")
		fmt.Printf("%s\n", cb)
	}

	return cnf
}

func isTrue(s string) bool {
	if s == "true" || s == "yes" || s == "1" {
		return true
	}
	return false
}

func compile(configs []config) (*Config, error) {
	cnf := &Config{}
	for _, c := range configs {

		cnf.Name = check(c.Name, cnf.Name).(string)
		cnf.DataCenter = check(c.DataCenter, cnf.DataCenter).(string)
		cnf.Host = check(c.Host, cnf.Host).(string)
		cnf.Port = check(c.Port, cnf.Port).(uint)
		cnf.IsServer = check(c.IsServer, cnf.IsServer).(bool)
		cnf.IsLeader = check(c.IsLeader, cnf.IsLeader).(bool)
		cnf.ServersAddr = check(c.ServersAddr, cnf.ServersAddr).(string)
		cnf.PrimaryServersAddr = check(c.PrimaryServersAddr, cnf.PrimaryServersAddr).(string)

		if c.Service != nil {
			cnf.Services = append(cnf.Services, Service(*c.Service))
		}
	}
	return cnf, nil
}

func loadConfigsFromDir(directory string) (cnf []config, err error) {
	var configs []config

	files, err := filepath.Glob(filepath.Join(directory, "*.hcl"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		var config config
		err = hclsimple.DecodeFile(file, nil, &config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func check(value any, default_value any) any {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return default_value
	}
	return v.Elem().Interface()
}
