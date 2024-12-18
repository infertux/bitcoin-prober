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
	./$(BINARY) --address 141.105.69.133
	./$(BINARY) --address bch.imaginary.cash:8333
	./$(BINARY) --address seed.bitnodes.io --network BTC

.PHONY: update-dependencies
update-dependencies:
	go get -u
	go mod tidy
