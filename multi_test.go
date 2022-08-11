package cache

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMultiCache(t *testing.T) {

	ctx := context.Background()

	dir, err := ioutil.TempDir("", "go-cache")

	if err != nil {
		t.Fatalf("Failed to create temp dir, %v", err)
	}

	defer os.RemoveAll(dir)

	cache_uri := fmt.Sprintf("multi://?cache=gocache://&cache=fs://%s", dir)

	c, err := NewCache(ctx, cache_uri)

	if err != nil {
		t.Fatalf("Failed to create multi:// cache, %v", err)
	}

	opts := &testCacheOptions{}

	// This is defined in testing.go
	err = testCache(ctx, c, opts)

	if err != nil {
		t.Fatalf("Cache tests failed, %v", err)
	}
}
