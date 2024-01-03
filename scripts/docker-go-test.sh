#!/bin/bash
LMD=/home/mrtdeh/go/pkg/mod
DMD=/go/pkg/mod
GOVER=1.20
docker run -it --rm -v $(pwd):/app -v $LMD:$DMD -w /app golang:$GOVER go test ./test/*