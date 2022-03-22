// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package manager

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/couchbaselabs/workbench-prototype/cluster-monitor/pkg/values"

	"github.com/couchbase/tools-common/restutil"
	"github.com/gorilla/mux"
)

type resultCluster struct {
	UUID           string                `json:"uuid"`
	Name           string                `json:"name"`
	NodesSummary   values.NodesSummary   `json:"nodes_summary"`
	BucketsSummary values.BucketsSummary `json:"buckets_summary"`
	HeartBeatIssue values.HeartIssue     `json:"heart_beat_issue,omitempty"`
	LastUpdate     time.Time             `json:"last_update"`
}

func (m *Manager) getClusterStatusReport(w http.ResponseWriter, r *http.Request) {
	m.getCheckerResultCommon(w, r, true)
}

func (m *Manager) getClusterStatusCheckerResult(w http.ResponseWriter, r *http.Request) {
	m.getCheckerResultCommon(w, r, false)
}

func (m *Manager) getCheckerResultCommon(w http.ResponseWriter, r *http.Request, filterDismissed bool) {
	vars := mux.Vars(r)
	uuid, ok := m.convertAliasToUUID(vars["uuid"], w)
	if !ok {
		return
	}

	cluster, err := m.store.GetCluster(uuid, false)
	if err != nil {
		if errors.Is(err, values.ErrNotFound) {
			restutil.HandleErrorWithExtras(restutil.ErrorResponse{
				Status: http.StatusNotFound,
				Msg:    fmt.Sprintf("cluster with UUID '%s' not found", uuid),
			}, w, nil)
			return
		}

		restutil.HandleErrorWithExtras(restutil.ErrorResponse{
			Status: http.StatusInternalServerError,
			Msg:    "could not get cluster details",
			Extras: err.Error(),
		}, w, nil)
		return
	}

	clusterOut := &resultCluster{
		UUID:           cluster.UUID,
		Name:           cluster.Name,
		BucketsSummary: cluster.BucketsSummary,
		NodesSummary:   cluster.NodesSummary,
		HeartBeatIssue: cluster.HeartBeatIssue,
		LastUpdate:     cluster.LastUpdate,
	}

	if err != nil {
		return
	}


	restutil.MarshalAndSend(http.StatusOK, clusterOut, w, nil)
}
