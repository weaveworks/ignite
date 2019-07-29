module github.com/weaveworks/ignite

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053 // indirect
	github.com/c2h5oh/datasize v0.0.0-20171227191756-4eba002a5eae
	github.com/containernetworking/cni v0.7.1
	github.com/containers/image v2.0.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/firecracker-microvm/firecracker-go-sdk v0.15.2-0.20190627223500-b2e8284e890c
	github.com/freddierice/go-losetup v0.0.0-20170407175016-fc9adea44124
	github.com/go-openapi/spec v0.17.0
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/gorilla/mux v1.7.2 // indirect
	github.com/krolaw/dhcp4 v0.0.0-20190531080455-7b64900047ae
	github.com/lithammer/dedent v1.1.0
	github.com/miekg/dns v1.1.14
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/weaveworks/flux v0.0.0-20190704153721-8292179855e1
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	google.golang.org/grpc v1.21.1 // indirect
	gopkg.in/alessio/shellescape.v1 v1.0.0-20170105083845-52074bc9df61
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/kube-openapi v0.0.0-20190401085232-94e1e7b7574c
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible
