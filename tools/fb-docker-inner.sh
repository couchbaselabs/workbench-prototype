#!/usr/bin/env bash
#
# Copyright (C) 2022 Couchbase, Inc.
#
# Use of this software is subject to the Couchbase Inc. License Agreement
# which may be found at https://www.couchbase.com/LA03012021.
#

set -exuo pipefail

# Runs as part of build-fluent-bit-docker.sh inside the container to build Fluent Bit.

if [ ! -f /.dockerenv ]; then
   echo "This script should only be run from inside a container."
   echo "If you're trying to build fluent-bit, use build-fluent-bit-docker.sh (or make fluent-bit-cross)."
   exit 1
fi

yum -y install epel-release # needed for cmake3
yum -y install bison cmake3 flex gcc gcc-c++ make openssl-devel

cd /work/upstream/fluent-bit/build

if [ -f Makefile ]; then
  make clean
  cd ..
  rm -rf build
  mkdir build
  cd build
fi
# disable in_systemd because systemd-libs isn't always present, dto for postgresql
cmake3 -DFLB_EXAMPLES=Off -DFLB_SHARED_LIB=Off -DBUILD_SHARED_LIBS=Off -DFLB_IN_SYSTEMD=Off -DFLB_OUT_PGSQL=Off ..
make "-j$(nproc)"
