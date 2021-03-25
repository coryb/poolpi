REMOTE_HOST = $(shell cat .poolpi-host)
export BUILDKIT_HOST ?= docker-container://buildkitd

push: console
	rsync -varH ./console $(REMOTE_HOST):~/poolpi/. 

console:
	hlb run -t console --log-output plain

lint:
	hlb run -t lint --log-output plain

protoc:
	hlb run -t protoc --log-output plain

################################################################################
# Setup targets
################################################################################

buildkitd:
	docker run -d --privileged --name buildkitd moby/buildkit:latest

hlb:
	go install github.com/openllb/hlb

.PHONY: push hlb buildkitd protoc lint console
