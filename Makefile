REMOTE_HOST = $(shell cat .poolpi-host)
export BUILDKIT_HOST ?= docker-container://buildkitd

push: console
	rsync -varH ./console $(REMOTE_HOST):~/poolpi/. 

console: console.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-s -w -extldflags "-static"' -o console *.go

lint:
	hlb run -t lint --log-output plain

protoc:
	 hlb run -t protoc --log-output plain

buildkitd:
	docker run -d --privileged --name buildkitd moby/buildkit:latest

hlb:
	go install github.com/openllb/hlb

.PHONY: push hlb buildkitd protoc lint
