package types

import (
	"time"
)

type AdMeta struct {
	Name       string
	At         uint8
	Aw         int
	Ah         int
	Size       int
	Priority   int
	AdToken    string
	Url        string
	UserId     int
	CreateTime time.Time
}
