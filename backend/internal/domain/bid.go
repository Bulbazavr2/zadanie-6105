package domain

import "time"

type Bid struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TenderID    string    `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorID    string    `json:"authorId"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

const (
	BidStatusPending  = "Created"
	BidStatusAccepted = "Published"
	BidStatusRejected = "Canceled"
)