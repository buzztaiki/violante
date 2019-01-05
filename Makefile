TARGETS = violante violante-server

all: $(TARGETS)

.PHONY: $(TARGETS)
$(TARGETS):
	go build ./cmd/$@
