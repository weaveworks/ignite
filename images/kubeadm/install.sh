#!/bin/bash

set -x

# Make kubectl commands work
export KUBECONFIG=/etc/kubernetes/admin.conf

MODE=${1}
BINARY_REF=${2}
CONTROL_PLANE_VERSION=${3}
DRY_RUN=${DRY_RUN:-0}
CLEANUP=${CLEANUP:-0}

if [[ $# != 3 && ${MODE} != "cleanup" ]]; then
	cat <<-EOF
	Usage:
		${0} MODE BINARY_REF CONTROL_PLANE_VERSION

	MODE=install|init|upgrade|cleanup
	BINARY_REF=What kubeadm binary and debs to pull. Can be a PR number, version label like "ci/latest" or exact (merged) commit like "v1.12.0-alpha.0-1035-gf2dec305ad"
	CONTROL_PLANE_VERSION=For init, this is the control plane version to use. For upgrade, this is the version to upgrade to
	EOF
	exit 1
fi

cleanup() {
	kubeadm reset -f
	if [[ -f $(which crictl) ]]; then
		# Should be called cri-tools
		apt-get purge cri-tools -y
	fi
	apt-get purge kubeadm kubelet kubernetes-cni kubectl -y
}

if [[ "${MODE}" == "cleanup" ]]; then
	cleanup
	exit 0
fi

if [[ ! -f $(which jq) ]]; then
	apt-get	update && apt-get install -y jq
fi
if [[ ! -f $(which docker) ]]; then
	apt-get	update && apt-get install -y docker.io
fi


BINARY_BUCKET=""
if [[ "${BINARY_REF}" =~ ^[0-9]{5}$ ]]; then
	PR_NUMBER=${BINARY_REF}
	BUILD_NUMBER=$(gsutil cat gs://kubernetes-jenkins/pr-logs/pull/${PR_NUMBER}/pull-kubernetes-bazel-build/latest-build.txt)
	BAZEL_PULL_REF=$(gsutil cat gs://kubernetes-jenkins/pr-logs/pull/${PR_NUMBER}/pull-kubernetes-bazel-build/${BUILD_NUMBER}/started.json | jq -r .pull)
	BAZEL_BUILD_LOCATION=$(gsutil cat gs://kubernetes-jenkins/shared-results/${BAZEL_PULL_REF}/bazel-build-location.txt)
	BINARY_BUCKET="${BAZEL_BUILD_LOCATION}/bin/linux/amd64"
elif [[ "${BINARY_REF}" =~ ^(ci|ci-cross){1}/latest ]]; then
	COMMIT=$(curl -sSL https://dl.k8s.io/${BINARY_REF}.txt)
	BINARY_BUCKET="gs://kubernetes-release-dev/ci/${COMMIT}-bazel/bin/linux/amd64"
elif [[ "${BINARY_REF}" =~ ^release/[a-z]+(-[0-9]+.[0-9]+)*$ ]]; then
	RELEASE=$(curl -sSL https://dl.k8s.io/${BINARY_REF}.txt)
	BINARY_BUCKET="gs://kubernetes-release/release/${RELEASE}/bin/linux/amd64"
else
	# Assume an exact "git describe" version/commit reference like "v1.12.0-alpha.0-1035-gf2dec305ad"
	BINARY_BUCKET="gs://kubernetes-release-dev/ci/${BINARY_REF}-bazel/bin/linux/amd64"
fi

# Download the debs and kubeadm
BINARY_DIR=$(mktemp -d)
gsutil cp ${BINARY_BUCKET}/*.deb ${BINARY_DIR}
gsutil cp ${BINARY_BUCKET}/kubeadm ${BINARY_DIR}
chmod +x ${BINARY_DIR}/kubeadm

install_kubeadm_apt() {
	apt-get update && apt-get install -y apt-transport-https curl || exit 0
	curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
	echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
	apt-get update && apt-get install -y kubeadm
}

install_debs() {
	# Install debs after upgrade
	dpkg -i ${BINARY_DIR}/*.deb
	apt-get update && apt-get install -f -y
	# bazel debs should ideally restart the kubelet automatically like k8s/release does
    # If we're only installing, remove the debs
	if [[ ${MODE} == "install" ]]; then
        rm -r ${BINARY_DIR}
    else
        systemctl start kubelet
    fi
}

if [[ "${MODE}" == "init" || "${MODE}" == "install" ]]; then
	INIT_KUBEADM=${BINARY_DIR}/kubeadm
	INIT_K8S_VERSION=${CONTROL_PLANE_VERSION}
	install_debs
elif [[ "${MODE}" == "upgrade" ]]; then
	INIT_KUBEADM="kubeadm"
	INIT_K8S_VERSION="stable"
	install_kubeadm_apt
fi

if [[ "${MODE}" == "install" ]]; then
    exit 0
fi

if [[ ${DRY_RUN} == 1 ]]; then
	${INIT_KUBEADM} init --dry-run --kubernetes-version ${INIT_K8S_VERSION}
fi

${INIT_KUBEADM} init --kubernetes-version ${INIT_K8S_VERSION} --ignore-preflight-errors=all
export KUBECONFIG=/etc/kubernetes/admin.conf
kubectl apply -f https://git.io/weave-kube-1.6

# Wait for the node to become ready
while [[ $(kubectl get node | tail -1 | awk '{print $2}') != "Ready" ]]; do sleep 1; done

if [[ "${MODE}" == "init" ]]; then
	exit 0
fi

${BINARY_DIR}/kubeadm upgrade plan --allow-experimental-upgrades

if [[ ${DRY_RUN} == 1 ]]; then
	${BINARY_DIR}/kubeadm upgrade apply --dry-run -f ${CONTROL_PLANE_VERSION}
fi

${BINARY_DIR}/kubeadm upgrade apply -f ${CONTROL_PLANE_VERSION}

# Install new debs
install_debs

kubeadm version
kubectl version
kubelet --version
crictl --version

if [[ ${CLEANUP} == 1 ]]; then
	cleanup
fi
