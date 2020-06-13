DC = docker-compose
CURRENT_DIR = $(shell pwd)

proto:
	docker run --rm -v $(CURRENT_DIR)/pb:/pb -v $(CURRENT_DIR)/schema:/proto ezio1119/protoc \
	-I/proto \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go:/pb" \
	event.proto chat.proto post.proto

up:
	${DC} up -d

logs:
	docker logs -f --tail 100 fishapp-relaylog_relaylog_1

down:
	${DC} down

build:
	${DC} build

clean:
	docker stop $(shell docker ps -aq)
	docker rm $(shell docker ps -aq)