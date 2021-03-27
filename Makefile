REMOTE_HOST = $(shell cat .poolpi-host)
export BUILDKIT_HOST ?= docker-container://buildkitd

push: binaries
	rsync -varH ./bin/ $(REMOTE_HOST):~/poolpi/. 

binaries:
	hlb run -t cmds --log-output plain

lint:
	hlb run -t lint --log-output plain

test:
	hlb run -t test --log-output plain

protoc:
	hlb run -t protoc --log-output plain

################################################################################
# Setup targets
################################################################################

buildkitd:
	docker run -d --privileged --name buildkitd moby/buildkit:latest

hlb:
	go install github.com/openllb/hlb

.PHONY: push hlb buildkitd protoc lint test binaries
