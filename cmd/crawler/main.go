package main

import (
	"flag"
	"guide-u/cafe-crawl/config"
	"guide-u/cafe-crawl/crawl"
	"guide-u/cafe-crawl/database"
	"log"
	"time"
)

func init() {
	config.LoadEnv()
}

func main() {
	// Flag
	startIndex, endIndex := 0, 0
	flag.IntVar(&startIndex, "start", 1, "collect start index")
	flag.IntVar(&endIndex, "end", 10, "collect end index")
	flag.Parse()

	start := time.Now()

	// Initial
	nc := crawl.NewNaverCrawl(
		config.Environment("NAVER_COOKIE_AUT"),
		config.Environment("NAVER_COOKIE_SES"),
	)

	err := database.DatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	go nc.Crawl()
	go nc.RunFailedArticle()

	// start crawl
	for i := startIndex; i <= endIndex; i++ {
		nc.Segement <- &crawl.ArticleSegement{
			Id:     int64(i),
			Status: crawl.Inactive,
		}
		time.Sleep(100 * time.Millisecond)
	}

	for {
		ok := true
		nc.Mutex.Lock()
		for _, v := range nc.Watcher {
			if !v {
				ok = false
			}
		}
		nc.Mutex.Unlock()

		if ok {
			close(nc.FailedSegement)
			close(nc.Segement)
			break
		}
	}

	end := time.Since(start)
	log.Println(end, nc.Count)
}
