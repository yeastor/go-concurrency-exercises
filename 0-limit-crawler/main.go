//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const rateLimit = time.Second // 10 calls per second
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("got interruption signal")
		cancel()
	}()

	crawler := NewCrawler()
	crawler.Crawl(ctx, rateLimit)
}
