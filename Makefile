BINARY := bitcoin-prober

all: $(BINARY)

$(BINARY): $(BINARY).go
	go build -v --race -o $(BINARY) $(BINARY).go
	strip $(BINARY)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: $(BINARY)
	./$(BINARY) --address 84.234.96.88
	./$(BINARY) --address seed.bchd.cash:8333
	./$(BINARY) --address seed.bitnodes.io --network BTC
