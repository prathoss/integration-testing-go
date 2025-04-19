package domain

import "time"

type Picture struct {
	ID        uint
	URL       string
	AuthorID  uint
	CreatedAt time.Time
	ViewCount uint
}
