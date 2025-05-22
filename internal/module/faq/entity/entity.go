package entity

import "time"

type FAQ struct {
	ID        int       `db:"id"`
	Question  string    `db:"question"`
	Answer    string    `db:"answer"`
	CreatedAt time.Time `db:"created_at"`
}
