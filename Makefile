DC = docker-compose
CURRENT_DIR = $(shell pwd)

proto:
	docker run --rm -v $(CURRENT_DIR)/src/domain:/work \
	-v $(CURRENT_DIR)/schema:/schema ezio1119/protoc \
	-I/schema \
	--doc_out=/schema \
	--doc_opt=markdown,README.md \
	--go_out=. \
	/schema/event/event.proto

up:
	${DC} up -d

logs:
	${DC} logs -f

down:
	${DC} stop nats-streaming
	${DC} rm nats-streaming

build:
	${DC} build