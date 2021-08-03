#!/bin/bash

set -o errexit
set -u

readonly commitHash=$(git log -n1 --format=format:'%H')
go build \
  -ldflags="-s -w -X github.com/99nil/arceus/pkg/version.version=$commitHash" \
  -installsuffix cgo \
  -o arceus \
  github.com/99nil/arceus/cmd