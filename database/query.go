package database

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// func AddPost(post Post) error {
// 	var existingPost Post
// 	isHave := connector.Where("articleid = ?", post.ArticleID).First(&existingPost)
// 	if isHave.Error != nil {
// 		if errors.Is(isHave.Error, gorm.ErrRecordNotFound) {
// 			// 레코드가 없으므로 생성
// 			tx := connector.Create(&post)
// 			return tx.Error
// 		} else {
// 			return errors.New("error add post")
// 		}
// 	} else {
// 		// 레코드가 있으므로 업데이트
// 		isHave.Updates(&post)
// 		if isHave.Error != nil {
// 			return errors.New("error update posts")
// 		}
// 	}

// 	return nil
// }

func APIAddPost(post Post) error {
	jsonData, err := json.Marshal(post)
	if err != nil {
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post("http://127.0.0.1:8000/article")

	fmt.Println(string(jsonData))

	if err != nil || resp.StatusCode() != 201 {
		return err
	}

	return nil
}
