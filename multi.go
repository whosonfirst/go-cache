package cache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"sync/atomic"
)

type MultiCache struct {
	Cache
	size      int64
	hits      int64
	misses    int64
	evictions int64
	caches    []Cache
}

func init() {
	ctx := context.Background()
	RegisterCache(ctx, "multi", NewMultiCache)
}

func NewMultiCache(ctx context.Context, str_uri string) (Cache, error) {

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	cache_uris := q["cache"]

	caches := make([]Cache, len(cache_uris))

	for idx, c_uri := range cache_uris {

		c, err := NewCache(ctx, c_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create cache for '%s', %v", c_uri, err)
		}

		caches[idx] = c
	}

	return NewMultiCacheWithCaches(ctx, caches...)
}

func NewMultiCacheWithCaches(ctx context.Context, caches ...Cache) (Cache, error) {

	c := &MultiCache{
		size:      int64(0),
		hits:      int64(0),
		misses:    int64(0),
		evictions: int64(0),
		caches:    caches,
	}

	return c, nil
}

func (mc *MultiCache) Close(ctx context.Context) error {

	for _, c := range mc.caches {

		err := c.Close(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

func (mc *MultiCache) Name() string {
	return "multi"
}

func (mc *MultiCache) Get(ctx context.Context, key string) (io.ReadCloser, error) {

	for _, c := range mc.caches {

		fh, err := c.Get(ctx, key)

		if err != nil {
			continue
		}

		atomic.AddInt64(&mc.hits, 1)

		// START OF...

		body, err := io.ReadAll(fh)

		if err != nil {
			return nil, err
		}

		go func(body []byte) {

			br := bytes.NewReader(body)

			for _, c := range mc.caches {

				// Only set caches that come *before* this cache
				
				if c.Name() == mc.Name() {
					break
				}

				cl := io.NopCloser(br)

				c.Set(ctx, key, cl)
				br.Seek(0, 0)
			}
		}(body)

		// END OF...

		br := bytes.NewReader(body)
		cl := io.NopCloser(br)

		return cl, nil
	}

	atomic.AddInt64(&mc.misses, 1)
	return nil, new(CacheMiss)
}

func (mc *MultiCache) Set(ctx context.Context, key string, fh io.ReadCloser) (io.ReadCloser, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)

	for _, c := range mc.caches {

		br.Seek(0, 0)
		cl := io.NopCloser(br)

		_, err := c.Set(ctx, key, cl)

		if err != nil {
			return nil, err
		}
	}

	atomic.AddInt64(&mc.size, 1)

	br.Seek(0, 0)
	cl := io.NopCloser(br)

	return cl, nil
}

func (mc *MultiCache) Unset(ctx context.Context, key string) error {

	for _, c := range mc.caches {

		err := c.Unset(ctx, key)

		if err != nil {
			return err
		}
	}

	atomic.AddInt64(&mc.size, -1)
	atomic.AddInt64(&mc.evictions, 1)

	return nil
}

func (mc *MultiCache) Size() int64 {
	return mc.SizeWithContext(context.Background())
}

func (mc *MultiCache) SizeWithContext(ctx context.Context) int64 {
	return atomic.LoadInt64(&mc.size)
}

func (mc *MultiCache) Hits() int64 {
	return atomic.LoadInt64(&mc.hits)
}

func (mc *MultiCache) Misses() int64 {
	return atomic.LoadInt64(&mc.misses)
}

func (mc *MultiCache) Evictions() int64 {
	return atomic.LoadInt64(&mc.evictions)
}
