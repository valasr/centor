package envoy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// NewEnvoyConfig creates a new EnvoyConfig instance
// with the given configuration settings and TLS configuration
//
// The configuration settings :
//
//	envoy.EnvoyConfig{
//		ListenerPort: 8081,
//		EndpointPort: 3000,
//		TLSConfig: envoy.TLSConfig{
//			Secure:         true,
//			CA:             "./pkg/envoy/testData/certs/ca.crt",
//			Cert:           "./pkg/envoy/testData/certs/server.crt",
//			Key:            "./pkg/envoy/testData/certs/server.key",
//			SessionTimeout: "6000s",
//		},
//	}
func NewEnvoy(cnf EnvoyConfig) error {
	strConfig := DefaultConfig

	if cnf.AccessLogPath == "" {
		cnf.AccessLogPath = "/var/log/centor-access.log"
	}
	if cnf.OutLogPath == "" {
		cnf.OutLogPath = "/var/log/centor-envoy.log"
	}
	if cnf.ListenerAddress == "" {
		cnf.ListenerAddress = "0.0.0.0"
	}
	if cnf.ListenerPort == 0 {
		cnf.ListenerPort = 80
	}
	if cnf.TLSConfig.Secure {
		strConfig = fmt.Sprintf(strConfig, downstreamTLS)
	} else {
		strConfig = fmt.Sprintf(strConfig, "")
	}

	if cnf.EndpointAddress == "" {
		cnf.EndpointAddress = "127.0.0.1"
	}
	if cnf.EndpointPort == 0 {
		return fmt.Errorf("endpoint port not specified")
	}

	strConfig = strings.ReplaceAll(strConfig, "{listener_address}", cnf.ListenerAddress)
	strConfig = strings.ReplaceAll(strConfig, "{listener_port}", fmt.Sprintf("%d", cnf.ListenerPort))
	strConfig = strings.ReplaceAll(strConfig, "{log_path}", cnf.AccessLogPath)
	strConfig = strings.ReplaceAll(strConfig, "{session_timeout}", cnf.TLSConfig.SessionTimeout)
	if cnf.TLSConfig.DisableSessionTicket {
		strConfig = strings.ReplaceAll(strConfig, "{disable_session_ticket}", "true")
	} else {
		strConfig = strings.ReplaceAll(strConfig, "{disable_session_ticket}", "false")
	}
	strConfig = strings.ReplaceAll(strConfig, "{ssl_cert}", cnf.TLSConfig.Cert)
	strConfig = strings.ReplaceAll(strConfig, "{ssl_key}", cnf.TLSConfig.Key)
	strConfig = strings.ReplaceAll(strConfig, "{ssl_ca}", cnf.TLSConfig.CA)
	strConfig = strings.ReplaceAll(strConfig, "{endpoint_address}", cnf.EndpointAddress)
	strConfig = strings.ReplaceAll(strConfig, "{endpoint_port}", fmt.Sprintf("%d", cnf.EndpointPort))

	// find the envoy binary path
	envoyBin, err := cnf.findBinary()
	if err != nil {
		return err
	}
	// envoy arguments
	args := []string{"--config-yaml", strConfig}
	// run the envoy process
	cmd := exec.Command(envoyBin, args...)
	if cnf.OutLogPath != "" {
		file, err := os.OpenFile(cnf.OutLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		cmd.Stdout = file
		cmd.Stderr = file
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (c *EnvoyConfig) findBinary() (string, error) {
	if c.EnvoyBinaryPath != "" {
		return c.EnvoyBinaryPath, nil
	}
	return exec.LookPath("envoy")
}
