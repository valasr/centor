#!/bin/bash
LMD=/home/mrtdeh/go/pkg/mod
DMD=/go/pkg/mod
GOVER=1.20
CONTAINER_NAME="docker-go-test"

if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}\$"; 
then

    echo "creating container \"$CONTAINER_NAME\""
    
    docker run -it \
    --name $CONTAINER_NAME \
    -v $(pwd):/app \
    -v $LMD:$DMD \
    -w /app golang:$GOVER sh -c "go clean -testcache && go test -v ./test/*"

else

    docker start $CONTAINER_NAME -i
fi





