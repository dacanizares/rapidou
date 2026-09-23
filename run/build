#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

mkdir -p bin
CGO_ENABLED=0 go build -trimpath -o bin/rapidou ./src
