#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

: "${APP_JWT_SECRET:=rapidou-local-secret-change-before-production}"
: "${APP_ADMIN_EMAIL:=admin@rapidou.test}"
: "${APP_ADMIN_PASSWORD:=rapidou-local-password}"
: "${APP_SEED_SAMPLE:=true}"

export APP_JWT_SECRET APP_ADMIN_EMAIL APP_ADMIN_PASSWORD APP_SEED_SAMPLE
exec go run ./src
