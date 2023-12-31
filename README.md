# Centor
Centor (Central Operations Management) is a distributed and highly available service implemented by the Raft protocol that provides a simple interface for securing, connecting, and monitoring nodes and services in clusters and can be used for different environments such as kubernetes or independent servers can be used.

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
## Deploy and Test
1. Build dockerFile :
```sh
make docker-build
```
2. Running up the clusters mesh :
```sh
make docker-run
```
3. Test cluster discovery with curl :
```sh
# You can call the following API's from this endpoints:
# :9991 => dc1, :9992 => dc2, :9993 => dc3, :9994 => dc4

# Ping all connected nodes in the current cluster and subclusters 
# and then return their names if available.
curl http://localhost:9991/call

# Display all connected nodes in the current cluster and sub-clusters
curl http://localhost:9991/nodes
```

