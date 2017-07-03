package main

import (
	"time"
)

type sleepLog struct {
	Sleep []struct {
		Date       string `json:"dateOfSleep"`
		DurationMS int64  `json:"duration"`
	} `json:"sleep"`
}

func (s sleepLog) MostRecent() string {
	log := s.Sleep[0]
	ms := time.Duration(log.DurationMS) * time.Millisecond
	return ms.String()
}
