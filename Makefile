TARGETS = violante violante-server

all: $(TARGETS)

.PHONY: $(TARGETS)
$(TARGETS):
	go build ./cmd/$@

.PHONY: rpc
rpc:
	protoc rpc/*.proto --go_out=plugins=grpc:.
