#!/usr/bin/env bash
set -euo pipefail

argocd_version="v3.5.2"

kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl apply \
  --namespace argocd \
  --server-side \
  --force-conflicts \
  --filename "https://raw.githubusercontent.com/argoproj/argo-cd/${argocd_version}/manifests/install.yaml"
kubectl wait \
  --namespace argocd \
  --for=condition=Available \
  deployment \
  --all \
  --timeout=300s
kubectl apply --filename k8s/argocd/root-application.yaml

echo "Argo CD ${argocd_version} is ready and the root application is installed."
