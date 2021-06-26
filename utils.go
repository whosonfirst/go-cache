package cache

import (
	"context"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"strings"
)

func SetString(c Cache, k string, v string) (string, error) {

	ctx := context.Background()

	r := strings.NewReader(v)
	rsc, err := ioutil.NewReadSeekCloser(r)

	if err != nil {
		return "", err
	}

	fh, err := c.Set(ctx, k, rsc)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	return toString(fh)
}

func GetString(c Cache, k string) (string, error) {

	ctx := context.Background()

	fh, err := c.Get(ctx, k)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	return toString(fh)
}

func toString(fh io.Reader) (string, error) {

	b, err := io.ReadAll(fh)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
