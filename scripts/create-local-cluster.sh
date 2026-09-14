#!/usr/bin/env bash
set -euo pipefail

cluster_name="planaria"
registry_name="planaria-registry"
registry_port="5001"

if ! docker inspect "${registry_name}" >/dev/null 2>&1; then
  docker run \
    --detach \
    --restart=always \
    --name "${registry_name}" \
    --publish "127.0.0.1:${registry_port}:5000" \
    registry:2.8.3
fi

if ! kind get clusters | grep -qx "${cluster_name}"; then
  kind create cluster --name "${cluster_name}" --config kind-cluster.yaml
fi

registry_directory="/etc/containerd/certs.d/localhost:${registry_port}"
for node in $(kind get nodes --name "${cluster_name}"); do
  docker exec "${node}" mkdir -p "${registry_directory}"
  printf '%s\n' \
    "[host.\"http://${registry_name}:5000\"]" \
    | docker exec -i "${node}" cp /dev/stdin "${registry_directory}/hosts.toml"
done

if [ "$(docker inspect -f '{{json .NetworkSettings.Networks.kind}}' "${registry_name}")" = "null" ]; then
  docker network connect kind "${registry_name}"
fi

kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${registry_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo "Kind cluster ${cluster_name} is ready with registry localhost:${registry_port}."
