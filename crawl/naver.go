package crawl

import (
	"encoding/json"
	"errors"
	"fmt"
	"guide-u/cafe-crawl/database"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// 게시물을 가져오는 URL (게시물 번호)
	N_POST_URL = `https://apis.naver.com/cafe-web/cafe-articleapi/v2.1/cafes/27842958/articles/%s`
	// 새 게시물을 가져오는 URL
	N_NEW_POST_URL = `https://apis.naver.com/cafe-web/cafe2/ArticleListV2dot1.json?search.clubid=27842958&search.queryType=lastArticle&search.page=%s&search.perPage=20`
	// 네이버 영상 정보을 가져오는 URL
	N_VIDEO_INFO_URL = `https://apis.naver.com/rmcnmv/rmcnmv/vod/play/v2.0/%s?key=%s`
	retry_count      = 3
)

var (
	errNoContent = errors.New("no content")
)

func NewNaverCrawl(aut, ses string) *NaverCrawl {
	client := resty.New()

	return &NaverCrawl{
		client: client,
		cookie: fmt.Sprintf("NID_AUT=%s;NID_SES=%s;",
			aut,
			ses,
		),
		Watcher:        make(map[int]bool, 0),
		Segement:       make(chan *ArticleSegement, 1024),
		FailedSegement: make(chan *FailedSegement, 1024),
		Mutex:          &sync.Mutex{},
		Count: &Count{
			Success: 0,
			Fail:    0,
		},
	}
}

func (nc *NaverCrawl) RunFailedArticle() {
	for v := range nc.FailedSegement {
		if v.Retry == 0 {
			log.Printf("재시도 횟수 초과 : %d\n", v.Id)
			nc.Count.Fail++
			nc.mutexWatcherWrite(int(v.Id), true)
			continue
		}

		if v.Status == Done {
			nc.Count.Success++
			nc.mutexWatcherWrite(int(v.Id), true)
			continue
		}

		log.Printf("%d 재시도, %d번 남음\n", v.Id, v.Retry)
		v.Retry -= 1
		nc.Segement <- &ArticleSegement{
			Id:     v.Id,
			Retry:  v.Retry,
			Status: Retry,
		}

		time.Sleep(1 * time.Second)
	}
}

func (nc *NaverCrawl) Crawl() {
	for v := range nc.Segement {
		if v.Status == Inactive {
			v.Status = Active
		}

		if v.Status != Retry {
			v.Retry = retry_count
		}

		article, err := nc.GetPost(int(v.Id))
		if err != nil {
			log.Printf("err : %v, id : %d\n", err, v.Id)
			if errors.Is(err, errNoContent) {
				nc.Count.Fail++
				nc.mutexWatcherWrite(int(v.Id), true)

				continue
			}

			nc.errorProcess(v)
			continue
		}

		dbArticle, err := nc.ParseArticle(v.Id, article)
		if err != nil {
			log.Printf("err : %v, id : %d\n", err, v.Id)
			nc.errorProcess(v)
			continue
		}

		//err = database.AddPost(dbArticle)
		err = database.APIAddPost(dbArticle)
		if err != nil {
			log.Printf("err : %v, id : %d\n", err, v.Id)
			nc.errorProcess(v)
			continue
		}

		log.Printf("DB 등록 완료 : %d\n", v.Id)
		v.Status = Done
		nc.Count.Success++
		nc.mutexWatcherWrite(int(v.Id), true)
	}
}

func (nc NaverCrawl) GetPost(num int) (Post, error) {
	url := fmt.Sprintf(N_POST_URL, strconv.Itoa(num))

	resp, err := nc.client.R().
		SetHeader("Cookie", nc.cookie).
		Get(url)

	if err != nil {
		log.Printf("error read post : %s\n", url)
		return Post{}, err
	}

	var post Post
	err = responseToMap(resp.Body(), &post)
	if err != nil {
		return post, err
	}

	if post.Result.ErrorCode != "" {
		switch post.Result.ErrorCode {
		case "4003":
			return post, errNoContent
		default:
			return post, fmt.Errorf("unknown error code : %v, reason : %v",
				post.Result.ErrorCode,
				post.Result.Reason)
		}
	}

	return post, nil
}

func (nc NaverCrawl) GetNewPost(page int) (newPost, error) {
	url := fmt.Sprintf(N_NEW_POST_URL, strconv.Itoa(page))

	resp, err := nc.client.R().
		SetHeader("Cookie", nc.cookie).
		Get(url)

	if err != nil {
		log.Println("error read post list")
	}

	var post newPost
	err = responseToMap(resp.Body(), &post)
	if err != nil {
		return post, err
	}

	if post.Message.Status != "200" {
		return post, fmt.Errorf("error read content : %v", post.Message.Error.Message)
	}

	return post, nil
}

func (nc NaverCrawl) GetVideoURL(vid, inKey string) (string, error) {
	url := fmt.Sprintf(N_VIDEO_INFO_URL, vid, inKey)

	resp, err := nc.client.R().
		SetHeader("Cookie", nc.cookie).
		Get(url)

	if err != nil {
		log.Printf("error read video : %s\n", url)
		return "", err
	}

	var video video
	err = responseToMap(resp.Body(), &video)
	if err != nil {
		return "", err
	}

	if len(video.Video.List) > 0 {
		// 최고화질
		return video.Video.List[len(video.Video.List)-1].Source, nil
	}

	return "", errors.New("no video")
}

func (nc NaverCrawl) errorProcess(v *ArticleSegement) {
	nc.FailedSegement <- &FailedSegement{
		Id:     v.Id,
		Retry:  v.Retry,
		Status: Retry,
	}

	nc.mutexWatcherWrite(int(v.Id), false)
}

func (nc NaverCrawl) mutexWatcherWrite(id int, value bool) {
	nc.Mutex.Lock()
	nc.Watcher[id] = value
	nc.Mutex.Unlock()
}

func responseToMap(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}
