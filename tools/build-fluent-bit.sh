#!/usr/bin/env bash

#
# Copyright (C) 2022 Couchbase, Inc.
#
# Use of this software is subject to the Couchbase Inc. License Agreement
# which may be found at https://www.couchbase.com/LA03012021.
#

set -euo pipefail

# Builds Fluent Bit from source. Note that this script builds it for your curent OS.
# If you want to build it for Linux using Docker, use build-fluent-bit-docker.sh.

FLB_REV=${FLB_REV:-master}

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [ -n "${OS:-}" ] && [ -n "${ARCH:-}" ]; then
  target="$SCRIPT_DIR/../build/fluent-bit-${OS}-${ARCH}"
else
  target="$SCRIPT_DIR/../build/fluent-bit"
fi

echo "Building Fluent Bit to $target"

pushd "$SCRIPT_DIR/../upstream/fluent-bit" || exit 1
  rm -rf build
  mkdir build
  pushd build || exit 1
    cmake -DFLB_EXAMPLES=Off -DFLB_SHARED_LIB=Off -DBUILD_SHARED_LIBS=OFF ..
    make "-j$(nproc)"
    cp bin/fluent-bit "$target"
  popd || exit 1
popd || exit 1
