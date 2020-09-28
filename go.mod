module github.com/weaveworks/ignite

go 1.14

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20190711223531-1fb7fffdb266
	github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
)

// TODO: Remove this when https://github.com/vishvananda/netlink/pull/554 is merged
replace github.com/vishvananda/netlink => github.com/twelho/netlink v1.1.1-ageing

require (
	github.com/alessio/shellescape v1.2.2
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/containerd/cgroups v0.0.0-20200407151229-7fc7a507c04c // indirect
	github.com/containerd/console v1.0.0
	github.com/containerd/containerd v1.3.3
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/containerd/go-cni v0.0.0-20200107172653-c154a49e2c75
	github.com/containerd/go-runc v0.0.0-20200220073739-7016d3ce2328 // indirect
	github.com/containerd/ttrpc v1.0.0 // indirect
	github.com/containerd/typeurl v1.0.0 // indirect
	github.com/containernetworking/plugins v0.8.5
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/go-iptables v0.4.5
	github.com/docker/docker v1.4.2-0.20200203170920-46ec8731fbce
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/firecracker-microvm/firecracker-go-sdk v0.21.1-0.20200312220944-e6eaff81c885
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/go-openapi/spec v0.19.8
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190909130307-a50d88189771
	github.com/lithammer/dedent v1.1.0
	github.com/miekg/dns v1.1.29
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nightlyone/lockfile v1.0.0
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runtime-spec v1.0.2
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.11.0
	github.com/prometheus/client_golang v1.5.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/vishvananda/netlink v1.1.0
	github.com/weaveworks/libgitops v0.0.0-20200611103311-2c871bbbbf0c
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4
	golang.org/x/tools v0.0.0-20200918232735-d647fc253266 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/apimachinery v0.19.2
	k8s.io/code-generator v0.19.2
	k8s.io/kube-openapi v0.0.0-20200831175022-64514a1d5d59
	sigs.k8s.io/yaml v1.2.0
)
