package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Crawler struct {
	throttle chan struct{}
}

func NewCrawler() *Crawler {
	return &Crawler{
		throttle: make(chan struct{}, 1),
	}
}

func (c *Crawler) Crawl(ctx context.Context, rateLimit time.Duration) {
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()
	defer fmt.Printf("cralw done.............. \n")

	var wg sync.WaitGroup
	go func() {
		for range ticker.C {
			select {
			case c.throttle <- struct{}{}:
			case <-ctx.Done():
				fmt.Printf("cancel ticker")
				return
			}
		}
	}()

	wg.Add(1)
	c.walk("http://golang.org/", 4, &wg, ctx)
	wg.Wait()
}

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func (c *Crawler) walk(url string, depth int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	if depth <= 0 {
		return
	}
	select {
	case <-ctx.Done():
		fmt.Printf("crawl done \n")
		return
	case <-c.throttle:

		fmt.Printf("tick, start crawl: %s \n", url)
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("found : %s %q\n", url, body)

		wg.Add(len(urls))
		for _, u := range urls {
			// Do not remove the `go` keyword, as Crawl() must be
			// called concurrently
			go c.walk(u, depth-1, wg, ctx)
		}
		return
	}

}
