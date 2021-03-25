module github.com/coryb/poolpi

go 1.16

require (
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/openllb/hlb v0.0.0-20210317202157-fe92ad094f28
	github.com/rs/xid v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
)

// all these are for openllb/hlb
replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/containerd/containerd => github.com/containerd/containerd v1.4.1-0.20201117152358-0edc412565dc

replace github.com/docker/cli => github.com/docker/cli v0.0.0-20200303162255-7d407207c304

replace github.com/docker/docker => github.com/docker/docker v20.10.0-beta1.0.20201110211921-af34b94a78a1+incompatible

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.0-pre1.0.20180209125602-c332b6f63c06

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20180806134042-1f13a808da65

replace github.com/sourcegraph/go-lsp => github.com/radeksimko/go-lsp v0.0.0-20200223162147-9f2c54f29c9f

replace github.com/tonistiigi/fsutil => github.com/slushie/fsutil v0.0.0-20200508061958-7d16a3dcbd1d
