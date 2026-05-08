#!/usr/bin/env bash
set -euo pipefail

export PATH="/Users/zhang/.nvm/versions/node/v22.20.0/bin:/opt/homebrew/bin:$PATH"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TAG="$(git -C "$ROOT" rev-parse --short HEAD)"
OWNER="${GHCR_OWNER:?set GHCR_OWNER}"
API_BASE_URL="${API_BASE_URL:?set API_BASE_URL}"
IMAGE="ghcr.io/${OWNER}/compx576-frontend:${TAG}"

docker buildx build --platform linux/amd64 -f "$ROOT/frontend.Dockerfile" -t "$IMAGE" --build-arg API_BASE_URL="$API_BASE_URL" --push "$ROOT"
kubectl -n compx576 set image deployment/frontend frontend="$IMAGE"
kubectl -n compx576 rollout status deployment/frontend
