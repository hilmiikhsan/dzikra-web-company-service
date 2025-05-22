package entity

import "time"

type ProductContent struct {
	ID          int       `db:"id"`
	ProductName string    `db:"product_name"`
	Images      string    `db:"images"`
	ContentID   string    `db:"content_id"`
	ContentEn   string    `db:"content_en"`
	SellLink    string    `db:"sell_link"`
	WebLink     string    `db:"web_link"`
	CreatedAt   time.Time `db:"created_at"`
}
