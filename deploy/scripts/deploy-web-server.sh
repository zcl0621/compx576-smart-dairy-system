#!/usr/bin/env bash
set -euo pipefail

export PATH="/Users/zhang/.nvm/versions/node/v22.20.0/bin:/opt/homebrew/bin:$PATH"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TAG="$(git -C "$ROOT" rev-parse --short HEAD)"
OWNER="${GHCR_OWNER:?set GHCR_OWNER}"
IMAGE="ghcr.io/${OWNER}/compx576-backend:${TAG}"

docker buildx build --platform linux/amd64 -f "$ROOT/backend.Dockerfile" -t "$IMAGE" --push "$ROOT"
kubectl -n compx576 set image deployment/web-server web-server="$IMAGE"
kubectl -n compx576 rollout status deployment/web-server
