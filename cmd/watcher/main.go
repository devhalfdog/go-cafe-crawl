package main

import (
	"guide-u/cafe-crawl/config"
	"guide-u/cafe-crawl/crawl"
	"guide-u/cafe-crawl/database"
	"log"
	"math"
	"time"
)

func init() {
	config.LoadEnv()
}

func main() {
	nc := crawl.NewNaverCrawl(
		config.Environment("NAVER_COOKIE_AUT"),
		config.Environment("NAVER_COOKIE_SES"),
	)

	err := database.DatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Initial Variable
	var (
		lastArticle int64 = 0
		watcher           = map[int]bool{}
	)

	go nc.Crawl()
	go nc.RunFailedArticle()

	for {
		page := 1
		for i := 1; i <= page; i++ {
			post, err := nc.GetNewPost(i)
			if err != nil {
				log.Println(err)
				break
			}

			if lastArticle != 0 {
				page = int(math.Ceil(float64(post.Message.Result.ArticleList[0].ID-lastArticle) / 20.0))
			}

			for _, v := range post.Message.Result.ArticleList {
				if v.ID > lastArticle {
					lastArticle = v.ID
				}

				go func(id int64) {
					if _, ok := watcher[int(id)]; !ok {
						nc.Segement <- &crawl.ArticleSegement{
							Id:     id,
							Status: crawl.Inactive,
						}
					}
				}(v.ID)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
