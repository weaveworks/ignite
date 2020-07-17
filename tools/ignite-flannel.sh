#!/bin/bash -e
# This is a helper script to set up a standalone CNI overlay using Flannel
# for use with Ignite VMs. See docs/networking.md for more information.

# Note: This script is meant to serve as a simple example, it does NOT
# secure the Flannel nor etcd traffic. Do not use it in production as-is.

shopt -s extglob nullglob

ETCD_VERSION=v3.2.27
FLANNEL_VERSION=v0.12.0-amd64
ETCD_NAME="ignite-etcd"
FLANNEL_NAME="ignite-flannel"

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

FLANNEL_NETWORK_CONFIG='{
	"Network": "10.50.0.0/16",
	"SubnetLen": 24,
	"SubnetMin": "10.50.10.0",
	"SubnetMax": "10.50.99.0",
	"Backend": {
		"Type": "udp",
		"Port": 8285
	}
}'

ARGS=("$@") # Save the args before shifting them away
ACTION_ARGS=() # Arguments passed to the action
ACTIONS="init|join|cleanup"

usage() {
	cat <<- EOF
	Usage:
	    $0 [-h|--help] ($ACTIONS) [parameter]...

	Available Actions:
	    init             Start etcd and Flannel on this machine. For standalone and multi-node host use cases.
	    join <host>      Start Flannel on this machine and join the multi-node overlay initialized by <host>.
	    cleanup          Remove all persistent data created by this script. Does not remove interfaces, iptables etc.

	Available Options:
	    -h, --help       Show this help text.
	EOF

	exit "$1"
}

join_usage() {
	cat <<- EOF
	Usage:
	    $0 join <host>

	where <host> is the IP address or FQDN of a host
	machine initialized with '$0 init'.
	EOF

	exit "$1"
}

# Set the given action
set_action() {
	[ -n "$ACTION" ] && usage 1
	ACTION="$1"
}

# Parse command line arguments
while test "$#" -gt 0; do
	case "$1" in
		@($ACTIONS))
			set_action "$1"
			;;
		-h|--help)
			usage 0
			;;
	esac

	shift

	[ -n "$ACTION" ] && ACTION_ARGS+=("$1")
done

[ -z "$ACTION" ] && usage 1

# Elevate privileges if needed
if [ "$EUID" -ne 0 ]; then
	exec sudo "$0" "${ARGS[@]}"
fi

# Fancy message formatting
log() {
	echo -e "\e[32;1m==>\e[39m $*\e[0m"
}

# Stops the given containers
stop_container() {
	mapfile -t ID <<< "$(docker ps -q --filter "name=$1")"
	{ [ -n "${ID[0]}" ] && docker stop "${ID[@]}" 1> /dev/null; } || true
}

# Forwards the passed parameters to 'etcdctl set'
set_config() {
	docker exec "$ETCD_NAME" etcdctl set "$@"
}

run_init() {
	log "Starting $ETCD_NAME container... "
	docker run -d --rm \
		-p 2379:2379 \
		-p 2380:2380 \
		--name "$ETCD_NAME" \
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

	log "Setting Flannel config:"
	set_config /coreos.com/network/config "$FLANNEL_NETWORK_CONFIG"

	log "Starting $FLANNEL_NAME container... "
	docker run -d --rm \
		--net host \
		--name "$FLANNEL_NAME" \
		--cap-add NET_ADMIN \
		--device /dev/net/tun \
		-v /etc/cni:/etc/cni \
		-v /run/flannel:/run/flannel \
		quay.io/coreos/flannel:${FLANNEL_VERSION} \
		-ip-masq

	log "Setting CNI config..."
	# Remove Ignite's conflist if it exists
	rm -f /etc/cni/net.d/10-ignite.conflist

	# Write the Flannel conflist
	mkdir -p /etc/cni/net.d
	echo "$FLANNEL_CONFLIST" > /etc/cni/net.d/10-flannel.conflist
}

run_join() {
	# Get the passed etcd endpoint
	[ -z "${ACTION_ARGS[0]}" ] && join_usage 1
	ETCD_ENDPOINT="http://${ACTION_ARGS[0]}:2379"

	log "Starting $FLANNEL_NAME container... "
	docker run -d --rm \
		--net host \
		--name "$FLANNEL_NAME" \
		--cap-add NET_ADMIN \
		--device /dev/net/tun \
		-v /etc/cni:/etc/cni \
		-v /run/flannel:/run/flannel \
		quay.io/coreos/flannel:${FLANNEL_VERSION} \
		-ip-masq \
		-etcd-endpoints "$ETCD_ENDPOINT"

	log "Setting CNI config..."
	# Remove Ignite's conflist if it exists
	rm -f /etc/cni/net.d/10-ignite.conflist

	# Write the Flannel conflist
	mkdir -p /etc/cni/net.d
	echo "$FLANNEL_CONFLIST" > /etc/cni/net.d/10-flannel.conflist
}

run_cleanup() {
	# Remove Flannel's conflist
	rm -f /etc/cni/net.d/10-flannel.conflist
}

# Stop etcd and/or Flannel if they're already running
stop_container "$ETCD_NAME|$FLANNEL_NAME"

# Run the specified action
case "$ACTION" in
	init)
		run_init
		log "Initialized, now start your Ignite VMs with the CNI network plugin."
		;;
	join)
		run_join
		log "Complete, now check if joining was successful using 'docker logs $FLANNEL_NAME'."
		log "If so, go ahead and start your Ignite VMs with the CNI network plugin."
		;;
	cleanup)
		run_cleanup
		log "Cleanup complete. To finish removal of non-persistent resources such as generated"
		log "network interfaces and iptables rules, reboot your system (or remove them by hand)."
		;;
	*)
		# Should never happen
		log "Internal error: invalid action $ACTION"
		exit 1
		;;
esac