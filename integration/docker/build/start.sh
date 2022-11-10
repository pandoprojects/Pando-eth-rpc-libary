#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

$GOBIN build -o ./build/linux/pando-eth-rpc-adaptor ./cmd/pando-eth-rpc-adaptor

set +x 

echo "Done."



