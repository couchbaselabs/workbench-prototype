// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package memcached

//go:generate mockery --name ConnIFace

// BucketCheckpointStats represents checkpoint statistics for a bucket, grouped by vBucket.
// The indexes of the slice elements are the indexes of the vBuckets.
type BucketCheckpointStats []map[string]string

// ConnIFace provides an interface to switch to a test memcached environment.
type ConnIFace interface {
	DCPStats(bucket string) ([]*DCPMemStats, error)
	MemStats(bucket string) ([]*MemoryStats, error)
	DefaultStats(bucket string) ([]*DefStats, error)
	CheckpointStats(host, bucket string) (BucketCheckpointStats, error)
	Hosts() []string
	Close() error
}
