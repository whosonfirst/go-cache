package cache

import (
	"context"
	"strings"
	"testing"
)

func TestGetSetString(t *testing.T) {

	k := "test"
	v := "test"

	ctx := context.Background()
	c, err := NewCache(ctx, "gocache://")

	if err != nil {
		t.Fatalf("Failed to create new cache, %v", err)
	}

	v2, err := SetString(ctx, c, k, v)

	if err != nil {
		t.Fatalf("Failed to set string, %v", err)
	}

	if strings.Compare(v, v2) != 0 {
		t.Fatalf("Strings don't match after setting: '%s', '%s'", v, v2)
	}

	v3, err := GetString(ctx, c, k)

	if err != nil {
		t.Fatalf("Failed to get string, %v", err)
	}

	if strings.Compare(v, v3) != 0 {
		t.Fatalf("Strings don't match after getting: '%s', '%s'", v, v3)
	}
}
