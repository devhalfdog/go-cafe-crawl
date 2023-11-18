package crawl

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

type NaverCrawl struct {
	client         *resty.Client
	cookie         string
	Watcher        map[int]bool
	Segement       chan *ArticleSegement
	FailedSegement chan *FailedSegement
	Mutex          *sync.Mutex
	Count          *Count
}

type Count struct {
	Success int
	Fail    int
}

type Status int

const (
	Active Status = iota
	Inactive
	Retry
	Done
)

type FailedSegement struct {
	Id     int64
	Retry  int
	Status Status
}

type ArticleSegement struct {
	Id     int64
	Retry  int
	Status Status
}
