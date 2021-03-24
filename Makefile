REMOTE_HOST=$(shell cat .poolpi-host)

push: console
	rsync -varH ./console $(REMOTE_HOST):~/poolpi/. 

console: console.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-s -w -extldflags "-static"' -o console *.go

.PHONY: push
