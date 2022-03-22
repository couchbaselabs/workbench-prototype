// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package manager

import (
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(m *Manager) *mux.Router {
	r := mux.NewRouter()

	r.Use(m.initializedMiddleware)
	r.Use(m.authMiddleware)
	r.Use(loggingMiddleware)

	metricsAPI(r)

	if m.config.EnableAdminAPI {
		adminAPI(r, m)
	}

	if m.config.EnableClusterAPI {
		clusterAPI(r, m)
	}

	if m.config.EnableExtendedAPI {
		extendedAPI(r, m)
	}

	if m.config.UIRoot != "" {
		ui(r, m)
	}

	return r
}

func adminAPI(r *mux.Router, m *Manager) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Administration and auth endpoints.
	// Get initialization state.
	v1.HandleFunc("/self", m.getInitState).Methods("GET")
	// Initialize cbmultimanager.
	v1.HandleFunc("/self", m.initializeCluster).Methods("POST")
	// Create JWT token endpoint.
	v1.HandleFunc("/self/token", m.tokenLogin).Methods("POST")

	zap.S().Info("(Routes) Set up Admin API")
}

func metricsAPI(r *mux.Router) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Collects prometheus metrics.
	v1.Handle("/_prometheus", promhttp.Handler()).Methods("GET")
	// Provide standard endpoint to simplify configuration.
	v1.Handle("/metrics", promhttp.Handler()).Methods("GET")

	zap.S().Info("(Routes) Set up Metrics API")
}

func clusterAPI(r *mux.Router, m *Manager) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Cluster management related endpoints.
	// Gets all the clusters.
	v1.HandleFunc("/clusters", m.getClusters).Methods("GET")
	// Adds a new cluster.
	v1.HandleFunc("/clusters", m.addNewCluster).Methods("POST")

	// Get only one specific cluster.
	v1.HandleFunc("/clusters/{uuid}", m.getCluster).Methods("GET")
	// Used to update the user, password, certificate or give a new bootstrap host.
	v1.HandleFunc("/clusters/{uuid}", m.updateClusterInfo).Methods("PATCH")
	// Stops tracking the cluster.
	v1.HandleFunc("/clusters/{uuid}", m.deleteCluster).Methods("DELETE")

	zap.S().Info("(Routes) Set up Cluster Management API")
}

func extendedAPI(r *mux.Router, m *Manager) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// UUID or bucket name using query parameters (bucket, node) respectively.
	v1.HandleFunc("/clusters/{uuid}/status/{name}", m.getClusterStatusCheckerResult).Methods("GET")

	// Get a single node's details (unblocker for https://issues.couchbase.com/browse/CMOS-188)
	v1.HandleFunc("/clusters/{uuid}/node/{node_uuid}", m.getClusterNodeDetails).Methods("GET")

	// Endpoints to manage cluster aliases.
	// Add alias endpoint.
	v1.HandleFunc("/aliases/{alias}", m.AddAlias).Methods("POST")
	// Delete alias endpoint.
	v1.HandleFunc("/aliases/{alias}", m.DeleteAlias).Methods("DELETE")

	// Endpoint to retrieve logs from the cluster.
	v1.HandleFunc("/clusters/{uuid}/nodes/{nodeUUID}/logs/{logName}", m.getLogs).Methods("GET")

	// Couchbase Cloud Endpoints
	v1.HandleFunc("/cloud/credentials", m.listCloudCreds).Methods("GET")
	v1.HandleFunc("/cloud/credentials", m.addCloudCreds).Methods("POST")

	v1.HandleFunc("/cloud/clusters", m.getCloudClusters).Methods("GET")
	v1.HandleFunc("/cloud/clusters/{id}", m.getCloudClusterStatus).Methods("GET")

	zap.S().Info("(Routes) Set up Extended API")
}

// ui serves the UI files under the UIRoot passed in the CLI. UI paths start with /ui. If the given path exists,
// is a file, and is not a hidden file, ui will serve it, otherwise it will serve index.html
// (with the assumption that the UI will handle the sub-path)
func ui(r *mux.Router, m *Manager) {
	r.PathPrefix(PathUIRoot).Methods("GET").Handler(
		http.StripPrefix(PathUIRoot, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPath := r.URL.Path
			if urlPath == "" || urlPath == "/" {
				http.ServeFile(w, r, path.Join(m.config.UIRoot, "index.html"))
				return
			}
			staticFile := path.Join(m.config.UIRoot, urlPath)
			stat, err := os.Stat(staticFile)
			if err == nil {
				if stat.IsDir() {
					w.WriteHeader(http.StatusForbidden)
				} else if base := path.Base(staticFile); base[0] == '.' {
					w.WriteHeader(http.StatusForbidden)
				} else {
					http.ServeFile(w, r, staticFile)
				}
			} else if errors.Is(err, os.ErrNotExist) {
				http.ServeFile(w, r, path.Join(m.config.UIRoot, "index.html"))
			} else {
				zap.S().Warnw("(Routes) Failed to serve static file", "err", err)
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Failed to serve UI file: " + err.Error()))
			}
		})))
	r.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(m.config.UIRoot, "index.html"))
	})

	zap.S().Infow("(Routes) Set up UI", "ui-root", m.config.UIRoot)
}
