REMOTE_HOST=$(shell cat .poolpi-host)

push: console
	rsync -varH ./console $(REMOTE_HOST):~/poolpi/. 

console: console.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-s -w -extldflags "-static"' -o console *.go

lint:
	docker run --rm -v $$(pwd):/src -w /src golangci/golangci-lint:v1.38.0 golangci-lint run ./...

.PHONY: push
