# Centor
----
Centor is a distributed, highly flexible, highly available service that provides a simple interface to secure, connect, and monitor nodes and services in the clusters and can be used for different environments such kubernetes or scrach server.

## Features
- **Multi-Datacenter** : easy to build multi datacenter clusters and create stable connection to gather.
- **K/V Database** : you can put your key/value pairs on local storage and access them in the other datacenters(dynamic app configuration).
- **System Health Check** : You can observe the health of the cluster's node.
- **Service Manager** : You can monitor and control services in the clusters.
- **Secure Bridge** : You can call an API or command on another node or cluster and get a response from it without worrying about the security of the connection.
- **Plaginable** : you can easy ro create awsome plugin to manage you'r environment in the all cluster.
- **Support HCL** : If you love HCL(hashicorp configuration language),ok leat's start that.

## Quick Start

You can simply build binary:
```sh
make build
```

And then print help for a more information :
```sh
./bin/centor -h
```
