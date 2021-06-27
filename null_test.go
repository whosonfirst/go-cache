package cache

import (
	"context"
	"testing"
)

func TestNullCache(t *testing.T) {

	ctx := context.Background()

	c, err := NewCache(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create null:// cache, %v", err)
	}

	opts := &testCacheOptions{
		AllowCacheMiss: true,
	}

	// This is defined in testing.go
	err = testCache(ctx, c, opts)

	if err != nil {
		t.Fatalf("Cache tests failed, %v", err)
	}
}
