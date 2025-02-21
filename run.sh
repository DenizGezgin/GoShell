#!/bin/sh

set -e
(
  cd "$(dirname "$0")"
  go build -o /tmp/shell-target cmd/shell/*.go
)
exec /tmp/shell-target