module github.com/weaveworks/ignite

go 1.16

replace github.com/docker/distribution => github.com/docker/distribution v0.0.0-20190711223531-1fb7fffdb266

// TODO: Remove this when https://github.com/vishvananda/netlink/pull/554 is merged
replace github.com/vishvananda/netlink => github.com/twelho/netlink v1.1.1-ageing

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/alessio/shellescape v1.2.2
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee
	github.com/containerd/cgroups v0.0.0-20210414185036-21be17332467 // indirect
	github.com/containerd/console v1.0.1
	github.com/containerd/containerd v1.5.0-beta.4
	github.com/containerd/continuity v0.0.0-20210417042358-bce1c3f9669b // indirect
	github.com/containerd/fifo v0.0.0-20210331061852-650e8a8a179d // indirect
	github.com/containerd/go-cni v1.0.1
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/containernetworking/plugins v0.8.7
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/go-iptables v0.4.5
	github.com/docker/cli v0.0.0-20200130152716-5d0cf8839492
	github.com/docker/docker v20.10.6+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/firecracker-microvm/firecracker-go-sdk v0.22.0
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/go-openapi/spec v0.19.8
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190909130307-a50d88189771
	github.com/lithammer/dedent v1.1.0
	github.com/miekg/dns v1.1.29
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nightlyone/lockfile v1.0.0
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runtime-spec v1.0.3-0.20200929063507-e6143ca7d51d
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.11.0
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/vishvananda/netlink v1.1.0
	github.com/weaveworks/libgitops v0.0.0-20200611103311-2c871bbbbf0c
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654
	golang.org/x/tools v0.1.8 // indirect
	google.golang.org/genproto v0.0.0-20210416161957-9910b6c460de // indirect
	google.golang.org/grpc v1.37.0 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/apimachinery v0.21.0
	k8s.io/code-generator v0.21.0
	k8s.io/gengo v0.0.0-20210203185629-de9496dff47b // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
	sigs.k8s.io/yaml v1.2.0
)
