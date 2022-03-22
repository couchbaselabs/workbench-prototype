// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package couchbase

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetAutoFailOverSettings() (*AutoFailoverSettings, error) {
	res, err := c.get(AutoFailOverSettings)
	if err != nil {
		return nil, fmt.Errorf("could not get auto failover settings: %w", err)
	}

	var settings AutoFailoverSettings
	if err = json.Unmarshal(res.Body, &settings); err != nil {
		return nil, fmt.Errorf("could not unmarshall the auto failover settings: %w", err)
	}

	return &settings, nil
}
