#!/usr/bin/env bash
set -euo pipefail

argocd_version="v3.5.2"

mise exec -- kubectl create namespace argocd --dry-run=client -o yaml \
  | mise exec -- kubectl apply -f -
mise exec -- kubectl apply \
  --namespace argocd \
  --server-side \
  --force-conflicts \
  --filename "https://raw.githubusercontent.com/argoproj/argo-cd/${argocd_version}/manifests/install.yaml"
mise exec -- kubectl wait \
  --namespace argocd \
  --for=condition=Available \
  deployment \
  --all \
  --timeout=300s
mise exec -- kubectl apply --filename k8s/argocd/root-application.yaml

echo "Argo CD ${argocd_version} is ready and the root application is installed."
