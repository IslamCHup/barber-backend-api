APP := ./cmd/barber

AIR := air


run:
	go run $(APP)

dev:
	$(AIR)

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
