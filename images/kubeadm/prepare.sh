#!/bin/bash

# Set up the seed node with the specified config file
IMAGE=${IMAGE:-"weaveworks/ignite-kubeadm"}

mkdir -p run
docker run -i --rm \
  -u $(id -u):$(id -g) \
  -v $(pwd)/run:/etc/kubernetes \
  ${IMAGE} \
    kubeadm init phase certs ca

docker run -i --rm \
  --net host \
  -u $(id -u):$(id -g) \
  -v $(pwd)/run:/etc/kubernetes \
  ${IMAGE} \
    kubeadm init phase kubeconfig admin

export HOST_IP=$(grep server run/admin.conf | grep -o -e "[0-9\.]*" | head -1)
export TOKEN=$(docker run -i --rm -v $(pwd)/run:/etc/kubernetes ${IMAGE} kubeadm token generate)
export CERT_KEY=$(docker run -i --rm -v $(pwd)/run:/etc/kubernetes ${IMAGE} kubeadm alpha certs certificate-key)
export CA_HASH=$(openssl x509 -pubkey -in run/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')

export LAST_ALLOCATED_IP=$(cat /var/lib/cni/networks/ignite-cni-bridge/last_reserved_ip.0)
export IP_PREFIX=$(echo ${LAST_ALLOCATED_IP} | cut -f-3 -d.)
export IP_START_NUMBER=$(echo ${LAST_ALLOCATED_IP} | cut -f4- -d.)

export MASTER1_IP="${IP_PREFIX}.$((IP_START_NUMBER + 1))"
export MASTER2_IP="${IP_PREFIX}.$((IP_START_NUMBER + 2))"
export MASTER3_IP="${IP_PREFIX}.$((IP_START_NUMBER + 3))"

cat > run/config.yaml <<EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
- token: "${TOKEN}"
certificateKey: "${CERT_KEY}"
nodeRegistration:
  criSocket: /run/containerd/containerd.sock
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: stable-1.18
controlPlaneEndpoint: ${MASTER1_IP}.xip.io:6443
apiServer:
  certSANs:
  - "${HOST_IP}"
EOF

cat > run/k8s-vars.sh <<EOF
export TOKEN=${TOKEN}
export CERT_KEY=${CERT_KEY}
export CA_HASH=${CA_HASH}
EOF

cat > run/haproxy.cfg <<EOF
frontend http_front
   bind *:443
   stats uri /haproxy?stats
   default_backend http_back

backend http_back
   balance roundrobin
   option httpchk GET /healthz
   http-check expect string ok
   server master1 ${MASTER1_IP}:6443 check check-ssl verify none
   server master2 ${MASTER2_IP}:6443 check check-ssl verify none
   server master3 ${MASTER3_IP}:6443 check check-ssl verify none
EOF
