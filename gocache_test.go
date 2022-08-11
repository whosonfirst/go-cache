package cache

import (
	"context"
	"testing"
)

func TestGoCache(t *testing.T) {

	ctx := context.Background()

	c, err := NewCache(ctx, "gocache://")

	if err != nil {
		t.Fatalf("Failed to create gocache cache, %v", err)
	}

	opts := &testCacheOptions{}

	// This is defined in testing.go
	err = testCache(ctx, c, opts)

	if err != nil {
		t.Fatalf("Cache tests failed, %v", err)
	}
}
