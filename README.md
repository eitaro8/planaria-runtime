# planaria-runtime

Local development runs on a Kind cluster backed by a Docker registry on
`localhost:5001`. The Kubernetes base is provider-neutral so that future GKE
overlays can replace only the registry, routing, scaling, and cloud-specific
configuration.

## Prerequisites

- Docker Desktop with WSL integration enabled for this distribution
- `mise install`

## Local environment

Create the Kind cluster and its local registry:

```sh
./scripts/create-local-cluster.sh
```

Build and publish the runtime image. The script publishes both a Git SHA tag
and the moving `dev` tag used only by the local overlay:

```sh
./scripts/publish-local-image.sh
```

After the changes containing the Argo CD configuration are merged into
`main`, bootstrap Argo CD and the root application:

```sh
./scripts/bootstrap-argocd.sh
```

Access the runtime without exposing it publicly:

```sh
kubectl --namespace planaria-local port-forward service/planaria-runtime 8080:80
curl --no-buffer http://localhost:8080/events
```

Access Grafana locally:

```sh
kubectl --namespace observability port-forward service/monitoring-grafana 3000:80
kubectl --namespace observability get secret monitoring-grafana \
  --output jsonpath='{.data.admin-password}' | base64 --decode
```

Open `http://localhost:3000` and sign in as `admin`. Prometheus and Tempo are
provisioned as data sources. No observability service is exposed outside the
cluster.
