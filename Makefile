DC = docker-compose
CURRENT_DIR = $(shell pwd)
API = post

sqldoc:
	docker run --rm --net=api-gateway_default -v $(CURRENT_DIR)/db:/work ezio1119/tbls \
	doc -f -t svg mysql://root:password@${API}-db:3306/${API}_DB ./

proto:
	docker run --rm -v $(CURRENT_DIR)/interfaces/controllers/${API}_grpc:$(CURRENT_DIR) \
	-v $(CURRENT_DIR)/schema/${API}:/schema \
	-w $(CURRENT_DIR) thethingsindustries/protoc \
	-I/schema \
	-I/usr/include/github.com/envoyproxy/protoc-gen-validate \
	--doc_out=markdown,README.md:/schema \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	${API}.proto

event:
	docker run --rm -v $(CURRENT_DIR)/interfaces/controllers/queue:$(CURRENT_DIR) \
	-v $(CURRENT_DIR)/schema/queue:/schema \
	-w $(CURRENT_DIR) thethingsindustries/protoc \
	-I/schema \
	-I/usr/include/github.com/envoyproxy/protoc-gen-validate \
	--doc_out=markdown,README.md:/schema \
	--go_out=plugins=grpc:. \
	--validate_out="lang=go:." \
	event.proto