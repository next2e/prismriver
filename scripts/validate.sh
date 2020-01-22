#!/bin/sh
# Heavily borrowed from the make scripts used by https://github.com/rancher.
set -e

PACKAGES="$(find . -name '*.go' | grep -Ev 'vendor|ab0x.go' | xargs -I{} dirname {} | sort -u)"
# Don't quote, we want this to glob.
go vet ${PACKAGES}
for i in ${PACKAGES}; do
  output="$(golint "$i")"
  if [ -n "$output" ]; then
    echo "$output"
    failed=true
  fi
done
test -z "$failed"
go fmt ${PACKAGES}
