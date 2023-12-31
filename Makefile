

.PHONY: all test clean

all: test build

test:
	go clean -testcache && go test -v ./test/*


docker-test:
	./scripts/docker-go-test.sh



build: clean
	go build -v -o bin/centor  ./main.go

protoc:
	protoc --go-grpc_out=require_unimplemented_servers=false:./proto/ ./proto/*.proto --go_out=./proto



docker-clean:
	docker image prune -f
	
docker-build: docker-clean
	docker build . --tag mrtdeh/centor
docker-up:
	docker compose -p dc1 -f ./docker-compose-dc1.yml up --force-recreate --build -d
	docker compose -p dc2 -f ./docker-compose-dc2.yml up --force-recreate -d
	docker compose -p dc3 -f ./docker-compose-dc3.yml up --force-recreate -d
	docker compose -p dc4 -f ./docker-compose-dc4.yml up --force-recreate -d

docker-down-all:
	docker compose -p dc1 -f ./docker-compose-dc1.yml down 
	docker compose -p dc2 -f ./docker-compose-dc2.yml down 
	docker compose -p dc3 -f ./docker-compose-dc3.yml down 
	docker compose -p dc4 -f ./docker-compose-dc4.yml down 
	docker compose -p dc1 -f ./docker-compose-with-envoy.yml down 




docker-build-with-envoy: docker-clean
	docker build -f dockerfile-with-envoy . --tag mrtdeh/centor:with-envoy
docker-up-with-envoy:
	docker compose -f docker-compose-with-envoy.yml -p dc1 up --force-recreate --build -d 



clean:
	rm -f ./bin/*