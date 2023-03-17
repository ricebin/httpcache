package httpcache_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bxcodec/httpcache"
	"github.com/redis/go-redis/v9"
)

func Example_inMemoryStorageDefault() {
	cached := httpcache.NewWithInmemoryCache(http.DefaultTransport, true, time.Second*15)
	client := &http.Client{
		Transport: cached,
	}

	processCachedRequest(client, cached)
	// Example Output:
	/*
		2020/06/21 13:14:51 Cache item's missing failed to retrieve from cache, trying with a live version
		Response time: 940086 micro-second
		Status Code 200
		Sequence >>>  0
		2020/06/21 13:14:53 Cache item's missing failed to retrieve from cache, trying with a live version
		Response time: 73679 micro-second
		Status Code 200
		Sequence >>>  1
		Response time: 126 micro-second
		Status Code 200
		Sequence >>>  2
		Response time: 96 micro-second
		Status Code 200
		Sequence >>>  3
		Response time: 102 micro-second
		Status Code 200
		Sequence >>>  4
		Response time: 94 micro-second
		Status Code 200
		Sequence >>>  5
	*/
}

func Example_redisStorage() {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer c.Close()

	cached := httpcache.NewWithRedisCache(http.DefaultTransport, true, c, time.Second*15)
	client := &http.Client{
		Transport: cached,
	}

	processCachedRequest(client, cached)
	// Example Output:
	/*
		2020/06/21 13:14:51 Cache item's missing failed to retrieve from cache, trying with a live version
		Response time: 940086 micro-second
		Status Code 200
		Sequence >>>  0
		2020/06/21 13:14:53 Cache item's missing failed to retrieve from cache, trying with a live version
		Response time: 73679 micro-second
		Status Code 200
		Sequence >>>  1
		Response time: 126 micro-second
		Status Code 200
		Sequence >>>  2
		Response time: 96 micro-second
		Status Code 200
		Sequence >>>  3
		Response time: 102 micro-second
		Status Code 200
		Sequence >>>  4
		Response time: 94 micro-second
		Status Code 200
		Sequence >>>  5
	*/
}

func processCachedRequest(client *http.Client, handler *httpcache.CacheHandler) {
	for i := 0; i < 100; i++ {
		startTime := time.Now()
		req, err := http.NewRequestWithContext(context.TODO(), "GET", "https://imantumorang.com", http.NoBody)
		if err != nil {
			log.Fatal((err))
		}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Response time: %v micro-second\n", time.Since(startTime).Microseconds())
		fmt.Println("Status Code", res.StatusCode)
		time.Sleep(time.Second * 1)
		fmt.Println("Sequence >>> ", i)
		if i%5 == 0 {
			err := handler.CacheInteractor.Flush(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
		res.Body.Close()
	}
}
