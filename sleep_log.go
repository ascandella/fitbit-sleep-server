package main

import (
	"time"
)

type sleep struct {
	Date       string `json:"dateOfSleep"`
	DurationMS int64  `json:"duration"`
}

type sleepLog struct {
	Sleep []sleep `json:"sleep"`
}

func (s sleep) MostRecent() string {
	ms := time.Duration(s.DurationMS) * time.Millisecond
	return ms.String()
}
