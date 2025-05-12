package data

import (
	"time"
)

type Grade struct {
	ID        int64
	CreatedAt time.Time
	Fullname  string
	Subject   string
	Grade     float64
	Email     string
}
