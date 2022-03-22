#!/usr/bin/env bash
#
# Copyright (C) 2022 Couchbase, Inc.
#
# Use of this software is subject to the Couchbase Inc. License Agreement
# which may be found at https://www.couchbase.com/LA03012021.
#

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

mkdir -p "$SCRIPT_DIR/../build/etc/fluent-bit"

cp -r "$SCRIPT_DIR/../upstream/couchbase-fluent-bit/conf/" "$SCRIPT_DIR/../build/etc/fluent-bit/"
rm "$SCRIPT_DIR"/../build/etc/fluent-bit/fluent-bit*.conf
cp -r "$SCRIPT_DIR"/../agent/pkg/fluentbit/conf/* "$SCRIPT_DIR/../build/etc/fluent-bit"
