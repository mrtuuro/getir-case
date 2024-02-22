package models

import (
	"time"
)

// Record struct helps us decode bson objects into Go structs.
type Record struct {
	Key        string    `bson:"key" json:"key"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	TotalCount int       `json:"totalCount"`
}
