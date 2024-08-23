package model

import "time"

type Audit struct {
	ID         int64
	Action     string
	CallParams string
	CreatedAt  time.Time
}
