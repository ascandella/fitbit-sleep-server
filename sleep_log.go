package main

import (
	"strings"
	"time"

	"go.uber.org/zap"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

type sleep struct {
	Date       string `json:"dateOfSleep"`
	DurationMS int64  `json:"duration"`
	StartTime  string `json:"startTime"`
}

type sleepLog struct {
	Sleep []sleep `json:"sleep"`
}

func (s sleep) MostRecent() string {
	ms := time.Duration(s.DurationMS) * time.Millisecond
	return ms.String()
}

func (s sleep) Start() string {
	chopped := strings.Split(s.StartTime, ".")[0]
	full := chopped
	t, err := time.Parse("2006-01-02T15:04:05", full)
	if err != nil {
		zap.L().Error("Couldn't parse time", zap.Error(err))
	}

	return t.In(location).Format(time.RFC1123)
}
