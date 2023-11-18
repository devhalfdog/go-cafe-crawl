package database

import "github.com/lib/pq"

type databaseConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	dbtype   string
}

type Post struct {
	ID              int64          `gorm:"column:id;auto_increment;primaryKey"`
	ArticleID       int64          `gorm:"column:articleid" json:"articleid"`
	Subject         string         `gorm:"column:subject" json:"subject"`
	OriginalContent string         `gorm:"column:original_content" json:"original_content"`
	Content         string         `gorm:"column:content" json:"content"`
	Image           pq.StringArray `gorm:"column:image;type:text[]" json:"image"`
	Video           pq.StringArray `gorm:"column:video;type:text[]" json:"video"`
	Menu            string         `gorm:"column:menu" json:"menu"`
	Head            string         `gorm:"column:head" json:"head"`
	Writer          string         `gorm:"column:writer" json:"writer"`
	CreateAt        int64          `gorm:"column:create" json:"create"`
}
