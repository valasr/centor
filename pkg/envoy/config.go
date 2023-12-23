package envoy

var (
	DefaultConfig = `{
		

		"static_resources": {
		  "listeners": [
			{
			  "address": {
				"socket_address": {
				  "address": "{listener_address}",
				  "port_value": {listener_port}
				}
			  },
			  "filter_chains": [
				{
				  "filters": [
					{
					  "name": "envoy.filters.network.http_connection_manager",
					  "typed_config": {
						"@type": "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
						"codec_type": "AUTO",
						"stat_prefix": "ingress_http",
						"access_log": [
						  {
							"name": "envoy.access_loggers.file",
							"typed_config": {
							  "@type": "type.googleapis.com/envoy.extensions.access_loggers.file.v3.FileAccessLog",
							  "log_format": {
								"json_format": {
								  "authority": "%%REQ(:AUTHORITY)%%",
								  "bytes_received": "%%BYTES_RECEIVED%%",
								  "bytes_sent": "%%BYTES_SENT%%",
								  "connection_termination_details": "%%CONNECTION_TERMINATION_DETAILS%%",
								  "downstream_local_address": "%%DOWNSTREAM_LOCAL_ADDRESS%%",
								  "downstream_remote_address": "%%DOWNSTREAM_REMOTE_ADDRESS%%",
								  "duration": "%%DURATION%%",
								  "method": "%%REQ(:METHOD)%%",
								  "path": "%%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%%",
								  "protocol": "%%PROTOCOL%%",
								  "request_id": "%%REQ(X-REQUEST-ID)%%",
								  "requested_server_name": "%%REQUESTED_SERVER_NAME%%",
								  "response_code": "%%RESPONSE_CODE%%",
								  "response_code_details": "%%RESPONSE_CODE_DETAILS%%",
								  "response_flags": "%%RESPONSE_FLAGS%%",
								  "route_name": "%%ROUTE_NAME%%",
								  "start_time": "%%START_TIME%%",
								  "upstream_cluster": "%%UPSTREAM_CLUSTER%%",
								  "upstream_host": "%%UPSTREAM_HOST%%",
								  "upstream_local_address": "%%UPSTREAM_LOCAL_ADDRESS%%",
								  "upstream_service_time": "%%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%%",
								  "upstream_transport_failure_reason": "%%UPSTREAM_TRANSPORT_FAILURE_REASON%%",
								  "user_agent": "%%REQ(USER-AGENT)%%",
								  "x_forwarded_for": "%%REQ(X-FORWARDED-FOR)%%"
								}
							  },
							  "path": "{log_path}"
							}
						  }
						],
						"route_config": {
						  "name": "local_route",
						  "virtual_hosts": [
							{
							  "name": "backend",
							  "domains": [
								"*"
							  ],
							  "routes": [
								{
								  "match": {
									"prefix": "/",
									"grpc": {
									}
								  },
								  "route": {
									"cluster": "backend_grpc_service"
								  }
								}
							  ]
							}
						  ]
						},
						"http_filters": [
						  {
							"name": "envoy.filters.http.router",
							"typed_config": {
							  "@type": "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router"
							}
						  }
						]
					  }
					}
				  ]
				  %s
				}
			  ]
			}
		  ],
		  "clusters": [
			{
			  "name": "backend_grpc_service",
			  "type": "STRICT_DNS",
			  "lb_policy": "ROUND_ROBIN",
			  "typed_extension_protocol_options": {
				"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": {
				  "@type": "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
				  "explicit_http_config": {
					"http2_protocol_options": {
					}
				  }
				}
			  },
			  "load_assignment": {
				"cluster_name": "backend_grpc_service",
				"endpoints": [
				  {
					"lb_endpoints": [
					  {
						"endpoint": {
						  "address": {
							"socket_address": {
							  "address": "{endpoint_address}",
							  "port_value": {endpoint_port}
							}
						  }
						}
					  }
					]
				  }
				]
			  }
			}
		  ]
		}
	  }`

	downstreamTLS = `
	,"transport_socket": {
		"name": "envoy.transport_sockets.tls",
		"typed_config": {
		  "@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext",
		  "session_timeout": "{session_timeout}",
		  "disable_stateless_session_resumption": {disable_session_ticket},
		  "common_tls_context": {
			"tls_certificates": [
			  {
				"certificate_chain": {
				  "filename": "{ssl_cert}"
				},
				"private_key": {
				  "filename": "{ssl_key}"
				}
			  }
			],
			"validation_context": {
			  "trusted_ca": {
				"filename": "{ssl_ca}"
			  }
			}
		  }
		}
	  }`
)

type TLSConfig struct {
	Secure               bool
	CA                   string
	Cert                 string
	Key                  string
	DisableSessionTicket bool
	SessionTimeout       string
}
type EnvoyConfig struct {
	EnvoyBinaryPath string

	OutLogPath      string
	AccessLogPath   string
	ListenerAddress string
	ListenerPort    uint
	TLSConfig
	EndpointAddress string
	EndpointPort    uint
}
