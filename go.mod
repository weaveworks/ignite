module github.com/weaveworks/ignite

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/Microsoft/hcsshim v0.8.6 // indirect
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053 // indirect
	github.com/c2h5oh/datasize v0.0.0-20171227191756-4eba002a5eae
	github.com/containerd/cgroups v0.0.0-20190717030353-c4b9ac5c7601 // indirect
	github.com/containerd/containerd v1.3.0-beta.1
	github.com/containerd/continuity v0.0.0-20190815185530-f2a389ac0a02 // indirect
	github.com/containerd/fifo v0.0.0-20190816180239-bda0ff6ed73c // indirect
	github.com/containerd/go-cni v0.0.0-20190813230227-49fbd9b210f3
	github.com/containerd/ttrpc v0.0.0-20190613183316-1fb3814edf44 // indirect
	github.com/containerd/typeurl v0.0.0-20190515163108-7312978f2987 // indirect
	github.com/containernetworking/cni v0.7.1 // indirect
	github.com/containers/image v2.0.0+incompatible
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/firecracker-microvm/firecracker-go-sdk v0.15.2-0.20190627223500-b2e8284e890c
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/fullsailor/pkcs7 v0.0.0-20190404230743-d7302db945fa // indirect
	github.com/go-openapi/spec v0.19.2
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/gogo/googleapis v1.2.0 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/gorilla/mux v1.7.2 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190531080455-7b64900047ae
	github.com/lithammer/dedent v1.1.0
	github.com/miekg/dns v1.1.14
	github.com/miscreant/miscreant-go v0.0.0-20190615163012-4f5dc8c406f6 // indirect
	github.com/miscreant/miscreant.go v0.0.0-20190615163012-4f5dc8c406f6 // indirect
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/otiai10/copy v1.0.1
	github.com/otiai10/curr v0.0.0-20190513014714-f5a3d24e5776 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/rjeczalik/notify v0.9.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f // indirect
	github.com/weaveworks/flux v0.0.0-20190704153721-8292179855e1
	github.com/whilp/git-urls v0.0.0-20160530060445-31bac0d230fa // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc // indirect
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/grpc v1.21.1 // indirect
	gopkg.in/alessio/shellescape.v1 v1.0.0-20170105083845-52074bc9df61
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/klog v0.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190709113604-33be087ad058
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible
	github.com/weaveworks/flux => ./third_party/forked/github.com/weaveworks/flux
)
