// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package couchbase

import (
	"github.com/couchbase/tools-common/cbrest"
)

const (
	PoolsNodesEndpoint       cbrest.Endpoint = "/pools/nodes"
	PoolsBucketEndpoint      cbrest.Endpoint = "/pools/default/buckets"
	PoolsBucketStatsEndpoint cbrest.Endpoint = "/pools/default/buckets/%s/stats"
	PoolsServerGroup         cbrest.Endpoint = "/pools/default/serverGroups"
	NodesSelfEndpoint        cbrest.Endpoint = "/nodes/self"

	UILogsEndpoint   cbrest.Endpoint = "/logs"
	SASLLogsEndpoint cbrest.Endpoint = "/sasl_logs/%s"

	AutoFailOverSettings cbrest.Endpoint = "/settings/autoFailover"

	PrometheusQueryEndpoint cbrest.Endpoint = "/_prometheus/api/v1/query_range"

	CheckersNodeEndpoint cbrest.Endpoint = "/_health/api/v1/checkers"
)
