package crawl

import (
	"encoding/json"
	"guide-u/cafe-crawl/database"
	"html"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (nc *NaverCrawl) ParseArticle(id int64, article Post) (database.Post, error) {
	dbArticle := ConvertPost(article)
	originalContent := html.UnescapeString(article.Result.Article.Content)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(originalContent))
	if err != nil {
		return database.Post{}, err
	}

	dbArticle.OriginalContent = originalContent
	dbArticle.Content = parseContent(originalContent, doc)
	dbArticle.Image = parseImage(originalContent, doc)
	dbArticle.Video, err = parseVideo(originalContent, doc, *nc)
	if err != nil {
		return database.Post{}, err
	}

	return dbArticle, nil
}

func ConvertPost(post Post) database.Post {
	return database.Post{
		ArticleID:       post.Result.Article.ID,
		Subject:         html.UnescapeString(post.Result.Article.Subject),
		Content:         "",
		OriginalContent: "",
		Image:           nil,
		Video:           nil,
		Menu:            html.UnescapeString(post.Result.Article.Menu.Name),
		Head:            html.UnescapeString(post.Result.Article.Head),
		Writer:          post.Result.Article.Writer.Nick,
		CreateAt:        post.Result.Article.WriteDate,
	}
}

func parseContent(body string, doc *goquery.Document) string {
	content := ""

	doc.Find(".se-text-paragraph").Each(func(i int, s *goquery.Selection) {
		content += s.Text() + "\n"
	})

	return content
}

func parseImage(body string, doc *goquery.Document) []string {
	content := []string{}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			content = append(content, src)
		}
	})

	return content
}

func parseVideo(body string, doc *goquery.Document, nc NaverCrawl) ([]string, error) {
	content := []string{}
	var err error

	// Youtube
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if strings.Contains(href, "youtu.be") || strings.Contains(href, "youtube.com") {
				content = append(content, href)
			}
		}
	})

	// Naver
	doc.Find("script.__se_module_data").Each(func(i int, s *goquery.Selection) {
		script, exists := s.Attr("data-module")
		if !exists {
			return
		}

		var data naverVideoData
		err = json.Unmarshal([]byte(script), &data)
		if err != nil {
			return
		}

		if data.Data.Inkey == "" || data.Data.Vid == "" {
			return
		}

		url, err := nc.GetVideoURL(data.Data.Vid, data.Data.Inkey)
		if err != nil {
			log.Println(err)
			return
		}

		if url != "" {
			content = append(content, url)
		}
	})

	return content, err
}
