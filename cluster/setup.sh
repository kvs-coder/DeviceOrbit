#!/bin/bash

# Ensure parameters are passed
if [[ $# -ne 2 ]]; then
  echo "Usage: $0 <IP_ADDRESS> <CLUSTER_NAME>"
  exit 1
fi

IP_ADDRESS=$1
CLUSTER_NAME=$2

sudo apt install -y docker.io k3d kubectx kubectl

function add_to_hosts {
    FILE='/etc/hosts'
    grep -xqF "$1" "$FILE" || echo "$1" >> "$FILE"
}

add_to_hosts "${IP_ADDRESS}    ${CLUSTER_NAME}"
add_to_hosts "${IP_ADDRESS}    ${CLUSTER_NAME}.local"

docker context rm -f "${CLUSTER_NAME}"
docker context create "${CLUSTER_NAME}" --description "Local ${CLUSTER_NAME} dev"
docker context use "${CLUSTER_NAME}"

sudo k3d cluster delete "${CLUSTER_NAME}"
sudo k3d cluster create -c k3d.yaml

k3d kubeconfig get $CLUSTER_NAME > $HOME/.kube/config
export KUBECONFIG="$HOME/.kube/config"