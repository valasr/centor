#!/bin/bash

docker rm -f $(docker ps -a|sed "1 d"| grep centor |awk '{print $1}'|sort)