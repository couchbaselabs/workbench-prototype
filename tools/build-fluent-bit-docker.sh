#!/usr/bin/env bash
#
# Copyright (C) 2022 Couchbase, Inc.
#
# Use of this software is subject to the Couchbase Inc. License Agreement
# which may be found at https://www.couchbase.com/LA03012021.
#

set -euo pipefail

# Builds Fluent Bit from source in a Docker container. Using centos7 to ensure
# we get an old enough (and thus widely compatible) glibc version.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Building Fluent Bit for $OS-$ARCH"

docker run \
    --rm -it \
    -v "$SCRIPT_DIR/..:/work" \
    centos:7 \
    /work/tools/fb-docker-inner.sh

mkdir -p "$SCRIPT_DIR/../build"
cp "$SCRIPT_DIR/../upstream/fluent-bit/build/bin/fluent-bit" "$SCRIPT_DIR/../build/fluent-bit-${OS}-${ARCH}"
