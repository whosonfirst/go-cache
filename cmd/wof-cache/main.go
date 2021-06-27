package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-cache"
	"log"
)

/*
go run cmd/wof-cache/main.go -cache-source 'gocache://' foo bar
2019/12/12 15:36:41 GET foo bar
2019/12/12 15:36:41 SET foo bar bar
2019/12/12 15:36:41 GET foo bar bar
*/

func main() {

	cache_source := flag.String("cache-source", "null://", "")

	flag.Parse()

	ctx := context.Background()

	c, err := cache.NewCache(ctx, *cache_source)

	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()

	if len(args)%2 == 1 {
		log.Fatal("Arguments not divisible by two (as in 'key -> value')")
	}

	for i := 0; i < len(args); i += 2 {

		k := args[i]
		v := args[i+1]

		g, err := cache.GetString(ctx, c, k)

		if err != nil && !cache.IsCacheMiss(err) {
			log.Fatal(err)
		}

		log.Println("GET", k, v, g)

		s, err := cache.SetString(ctx, c, k, v)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("SET", k, v, s)

		v2, err := cache.GetString(ctx, c, k)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("GET", k, v, v2)
	}

}
