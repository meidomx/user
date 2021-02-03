package model

import (
	"errors"
	"time"
)

var ErrNoRecordFound = errors.New("no record found")

func UnixTime(t time.Time) int64 {
	return t.Unix()*1000 + t.UnixNano()/1000/1000%1000
}
