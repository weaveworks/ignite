#!/bin/bash

# Set up the seed node with the specified config file
mkdir -p run
docker run -i --rm -v $(pwd)/run:/etc/kubernetes weaveworks/ignite-kubeadm \
    kubeadm init phase certs ca

docker run -i --rm --net host -v $(pwd)/run:/etc/kubernetes weaveworks/ignite-kubeadm \
    kubeadm init phase kubeconfig admin

export HOST_IP=$(grep server run/admin.conf | grep -o -e "[0-9\.]*" | head -1)
export TOKEN=$(docker run -i --rm -v $(pwd)/run:/etc/kubernetes weaveworks/ignite-kubeadm kubeadm token generate)
export CERT_KEY=$(docker run -i --rm -v $(pwd)/run:/etc/kubernetes weaveworks/ignite-kubeadm kubeadm alpha certs certificate-key)
export CA_HASH=$(openssl x509 -pubkey -in run/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')

cat > run/config.yaml <<EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
- token: "${TOKEN}"
certificateKey: "${CERT_KEY}"
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v1.15.0-rc.1
controlPlaneEndpoint: firekube.luxas.dev:6443
apiServer:
  certSANs:
  - "${HOST_IP}"
EOF

cat > run/k8s-vars.sh <<EOF
export TOKEN=${TOKEN}
export CERT_KEY=${CERT_KEY}
export CA_HASH=${CA_HASH}
EOF
