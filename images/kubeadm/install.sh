#!/bin/bash

set -x

# Make kubectl commands work
export KUBECONFIG=/etc/kubernetes/admin.conf

MODE=${1}
BINARY_REF=${2}
CONTROL_PLANE_VERSION=${3}
GOARCH=${4}
DRY_RUN=${DRY_RUN:-0}
CLEANUP=${CLEANUP:-0}

if [[ $# != 4 && ${MODE} != "cleanup" ]]; then
	cat <<-EOF
	Usage:
		${0} MODE BINARY_REF CONTROL_PLANE_VERSION

	MODE=install|init|upgrade|cleanup
	BINARY_REF=What kubeadm binary and debs to pull. Can be a PR number, version label like "ci/latest" or exact (merged) commit like "v1.12.0-alpha.0-1035-gf2dec305ad"
	CONTROL_PLANE_VERSION=For init, this is the control plane version to use. For upgrade, this is the version to upgrade to
	GOARCH=amd64|arm|arm64|ppc64le|s390x
	EOF
	exit 1
fi

gsutil_cp() {
	if [[ -f $(which gsutil 2>/dev/null) ]]; then
		gsutil cp $1 $2
	else
		URL=$(echo $1 | sed "s|gs://|https://storage.googleapis.com/|")
		curl -sSL ${URL} > $2
	fi
}

gsutil_cat() {
	if [[ -f $(which gsutil 2>/dev/null) ]]; then
		gsutil cat $1
	else
		URL=$(echo $1 | sed "s|gs://|https://storage.googleapis.com/|")
		curl -sSL ${URL}
	fi
}

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

BINARY_BUCKET=""
if [[ "${BINARY_REF}" =~ ^[0-9]{5}$ ]]; then
	PR_NUMBER=${BINARY_REF}
	BUILD_NUMBER=$(gsutil_cat gs://kubernetes-jenkins/pr-logs/pull/${PR_NUMBER}/pull-kubernetes-bazel-build/latest-build.txt)
	BAZEL_PULL_REF=$(gsutil_cat gs://kubernetes-jenkins/pr-logs/pull/${PR_NUMBER}/pull-kubernetes-bazel-build/${BUILD_NUMBER}/started.json | jq -r .pull)
	BAZEL_BUILD_LOCATION=$(gsutil_cat gs://kubernetes-jenkins/shared-results/${BAZEL_PULL_REF}/bazel-build-location.txt)
	BINARY_BUCKET="${BAZEL_BUILD_LOCATION}/bin/linux/${GOARCH}"
elif [[ "${BINARY_REF}" =~ ^(ci|ci-cross){1}/latest ]]; then
	COMMIT=$(curl -sSL https://dl.k8s.io/${BINARY_REF}.txt)
	BINARY_BUCKET="gs://kubernetes-release-dev/ci/${COMMIT}-bazel/bin/linux/${GOARCH}"
elif [[ "${BINARY_REF}" =~ ^release/[a-z]+(-[0-9]+.[0-9]+)*$ ]]; then
	RELEASE=$(curl -sSL https://dl.k8s.io/${BINARY_REF}.txt)
	BINARY_BUCKET="gs://kubernetes-release/release/${RELEASE}/bin/linux/${GOARCH}"
	INSTALL_APT=true
else
	# Assume an exact "git describe" version/commit reference like "v1.12.0-alpha.0-1035-gf2dec305ad"
	BINARY_BUCKET="gs://kubernetes-release-dev/ci/${BINARY_REF}-bazel/bin/linux/${GOARCH}"
fi

# Download the debs and kubeadm
BINARY_DIR=$(mktemp -d)
for pkg in cri-tools kubeadm kubectl kubelet kubernetes-cni; do
	echo "Downloading ${pkg}.deb"
	gsutil_cp ${BINARY_BUCKET}/${pkg}.deb ${BINARY_DIR}/${pkg}.deb
done

gsutil_cp ${BINARY_BUCKET}/kubeadm ${BINARY_DIR}/kubeadm
chmod +x ${BINARY_DIR}/kubeadm

install_kubeadm_apt() {
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
	if [[ ${INSTALL_APT} == "true" ]]; then
		install_kubeadm_apt
	else
		install_debs
	fi
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
 