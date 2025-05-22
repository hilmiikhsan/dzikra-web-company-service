package entity

import "time"

type Article struct {
	ID                  int       `db:"id"`
	Title               string    `db:"title"`
	Image               string    `db:"image"`
	Content             string    `db:"content"`
	Description         string    `db:"description"`
	ArticleCategoryID   int       `db:"article_category_id"`
	ArticleCategoryName string    `db:"article_category_name"`
	CreatedAt           time.Time `db:"created_at"`
}
