#!/usr/bin/env bash
set -euo pipefail

repository="localhost:5001/planaria-runtime"
git_tag="$(git rev-parse --short=12 HEAD)"

docker build \
  --tag "${repository}:${git_tag}" \
  --tag "${repository}:dev" \
  .

docker push "${repository}:${git_tag}"
docker push "${repository}:dev"

echo "Published ${repository}:${git_tag} and updated the local dev tag."
