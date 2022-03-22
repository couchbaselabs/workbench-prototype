// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package discovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/couchbaselabs/workbench-prototype/cluster-monitor/pkg/discovery/mocks"
)

func TestDiscoveryLoop(t *testing.T) {
	t.Parallel()
	mockDisco := mocks.CouchbaseClusterDiscovery{}
	mockDisco.On("Discover", mock.Anything).Return(nil)
	dm, err := NewDiscoveryManager(&mockDisco)
	require.NoError(t, err)
	dm.Start(time.Second)
	time.Sleep(1500 * time.Millisecond)
	mockDisco.AssertNumberOfCalls(t, "Discover", 2)
	dm.Stop()
	mockDisco.AssertNumberOfCalls(t, "Discover", 2)
}

func TestDiscoveryLoopError(t *testing.T) {
	t.Parallel()
	mockDisco := mocks.CouchbaseClusterDiscovery{}
	mockDisco.On("Discover", mock.Anything).Return(fmt.Errorf("oh no"))
	dm, err := NewDiscoveryManager(&mockDisco)
	require.NoError(t, err)
	dm.Start(time.Second)
	time.Sleep(1500 * time.Millisecond)
	mockDisco.AssertNumberOfCalls(t, "Discover", 2)
	dm.Stop()
	mockDisco.AssertNumberOfCalls(t, "Discover", 2)
}

func TestDiscoveryStopBeforeStart(t *testing.T) {
	mockDisco := mocks.CouchbaseClusterDiscovery{}
	mockDisco.On("Discover", mock.Anything).Return(nil)
	dm, err := NewDiscoveryManager(&mockDisco)
	require.NoError(t, err)
	dm.Stop()
	mockDisco.AssertNumberOfCalls(t, "Discover", 0)
}
