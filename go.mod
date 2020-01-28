module github.com/weaveworks/ignite

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053 // indirect
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/containerd/cgroups v0.0.0-20200116170754-a8908713319d // indirect
	github.com/containerd/console v0.0.0-20191219165238-8375c3424e4d
	github.com/containerd/containerd v1.3.2
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41 // indirect
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00 // indirect
	github.com/containerd/go-cni v0.0.0-20200107172653-c154a49e2c75
	github.com/containerd/go-runc v0.0.0-20191213151417-8aa1fd617053 // indirect
	github.com/containerd/ttrpc v0.0.0-20200121165050-0be804eadb15 // indirect
	github.com/containerd/typeurl v0.0.0-20200115183213-fe1d0d650e42 // indirect
	github.com/containernetworking/plugins v0.8.5
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/go-iptables v0.4.5
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20191113042239-ea84732a7725
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/firecracker-microvm/firecracker-go-sdk v0.15.2-0.20190627223500-b2e8284e890c
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/go-openapi/spec v0.19.5
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/kr/pty v1.1.5 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190909130307-a50d88189771
	github.com/lithammer/dedent v1.1.0
	github.com/miekg/dns v1.1.27
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.1
	github.com/otiai10/copy v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/vishvananda/netlink v1.1.0
	github.com/weaveworks/gitops-toolkit v0.0.0-20191009213125-6d6e6c8de6d4
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/alessio/shellescape.v1 v1.0.0-20170105083845-52074bc9df61
	gotest.tools v2.2.0+incompatible
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
	k8s.io/apimachinery v0.17.2
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/docker/distribution => github.com/docker/distribution v0.0.0-20190711223531-1fb7fffdb266
