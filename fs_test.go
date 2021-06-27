package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFSCache(t *testing.T) {

	ctx := context.Background()

	dir, err := ioutil.TempDir("", "go-cache")

	if err != nil {
		t.Fatalf("Failed to create temp dir, %v", err)
	}

	defer os.RemoveAll(dir)

	cache_uri := fmt.Sprintf("fs://%s", dir)

	c, err := NewCache(ctx, cache_uri)

	if err != nil {
		t.Fatalf("Failed to create fs:// cache, %v", err)
	}

	opts := &testCacheOptions{}

	// This is defined in testing.go
	err = testCache(ctx, c, opts)

	if err != nil {
		t.Fatalf("Cache tests failed, %v", err)
	}
}
