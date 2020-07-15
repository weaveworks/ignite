#!/bin/bash -e

ETCD_VERSION=v3.2.27
FLANNEL_VERSION=v0.12.0-amd64
DOCKER="sudo docker"

FLANNEL_CONFLIST='{
	"name": "cbr0",
	"cniVersion": "0.4.0",
	"plugins": [
		{
			"type": "flannel",
			"delegate": {
				"hairpinMode": true,
				"isDefaultGateway": true
			}
		},
		{
			"type": "portmap",
			"capabilities": {
				"portMappings": true
			}
		},
		{
			"type": "firewall"
		}
	]
}'

root_write() {
	echo "$2" | sudo tee "$1" > /dev/null
}

stop_container() {
	mapfile -t ID <<< "$(${DOCKER} ps -q --filter "name=$1")"
	{ (( ${#ID[@]} )) && ${DOCKER} rm -f "${ID[@]}" 1> /dev/null; } || true
}

start_etcd() {
	${DOCKER} run -d --rm \
		-p 2379:2379 \
		-p 2380:2380 \
		--name "$1" \
		gcr.io/etcd-development/etcd:${ETCD_VERSION} \
 		/usr/local/bin/etcd \
		--name s1 \
		--data-dir /etcd-data \
		--listen-client-urls http://0.0.0.0:2379 \
		--advertise-client-urls http://0.0.0.0:2379 \
		--listen-peer-urls http://0.0.0.0:2380 \
		--initial-advertise-peer-urls http://0.0.0.0:2380 \
		--initial-cluster s1=http://0.0.0.0:2380 \
		--initial-cluster-token tkn \
		--initial-cluster-state new
}

start_flannel() {
	${DOCKER} run -d --rm \
		--net host \
		--name "$1" \
		--cap-add NET_ADMIN \
		--device /dev/net/tun \
		-v /etc/cni:/etc/cni \
		-v /run/flannel:/run/flannel \
		quay.io/coreos/flannel:${FLANNEL_VERSION} \
		-ip-masq
}

set_config() {
	# Forward the passed parameters to 'etcdctl set'
	${DOCKER} exec ignite-etcd etcdctl set "$@"
}

set_cni_conf() {
	# Remove Ignite's conflist if it exists
	rm -f /etc/cni/net.d/10-ignite.conflist

	# Write the Flannel conflist
	root_write /etc/cni/net.d/10-flannel.conflist "$FLANNEL_CONFLIST"
}

# Stop etcd and/or Flannel if they're already running
stop_container "ignite-etcd|ignite-flannel"

echo "Starting ignite-etcd container:"
start_etcd ignite-etcd

echo "Setting Flannel config:"
set_config /coreos.com/network/config '{ "Network": "10.1.0.0/16" }'

echo "Starting ignite-flannel container:"
start_flannel ignite-flannel

echo "Setting CNI config."
set_cni_conf

echo "Initialized, now run your VMs with CNI networking."