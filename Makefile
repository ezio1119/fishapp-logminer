CWD = $(shell pwd)
SVC = relaylog
CHAT_DB_SVC = chat-db
POST_DB_SVC = post-db
NATS_URL = nats-streaming:4223
NET = fishapp-net
PJT_NAME = $(notdir $(PWD))
# TEST = $(shell docker inspect $(NET) > /dev/null 2>&1; echo " $$?")

createnet:
	docker network create $(NET)

proto:
	docker run --rm --name protoc -v $(CWD)/pb:/pb -v $(CWD)/schema:/proto ezio1119/protoc \
	-I/proto \
	-I/go/src/github.com/envoyproxy/protoc-gen-validate \
	--go_out=plugins=grpc:/pb \
	--validate_out="lang=go:/pb" \
	event.proto chat.proto post.proto

waitchatdb:
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-timeout 30s \
	-wait tcp://$(CHAT_DB_SVC):3306

waitpostdb:
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-timeout 30s \
	-wait tcp://$(POST_DB_SVC):3306

waitnats:
	docker run --rm --name dockerize --net $(NET) jwilder/dockerize \
	-wait tcp://$(NATS_URL)

test:
	docker-compose exec $(SVC) sh -c "go test -v -coverprofile=cover.out ./... && \
	go tool cover -html=cover.out -o ./cover.html" && \
	open ./src/cover.html

up: waitchatdb waitpostdb waitnats
	docker-compose up -d $(SVC)

build:
	docker-compose build

down:
	docker-compose down

exec:
	docker-compose exec $(SVC) sh

logs:
	docker logs -f --tail 100 $(PJT_NAME)_$(SVC)_1

rmvol:
	docker-compose down -v