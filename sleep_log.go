package main

import (
	"bytes"
	"fmt"
	"math"
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
	Date          string `json:"dateOfSleep"`
	MinutesAsleep int    `json:"minutesAsleep"`
	StartTime     string `json:"startTime"`
}

type sleepLog struct {
	Sleep []sleep `json:"sleep"`
}

func (s sleep) FriendlyDuration() string {
	ms := time.Duration(s.MinutesAsleep) * time.Minute
	out := &bytes.Buffer{}
	if ms.Hours() >= 1.0 {
		fmt.Fprintf(out, "%.0f hour", ms.Hours())
		if ms.Hours() >= 2.0 {
			fmt.Fprint(out, "s")
		}
		fmt.Fprint(out, ", ")
	}
	fmt.Fprintf(out, "%.0f minutes", math.Mod(ms.Minutes(), 60.0))
	return out.String()
}

func (s sleep) Start() string {
	chopped := strings.Split(s.StartTime, ".")[0]
	full := chopped
	t, err := time.ParseInLocation("2006-01-02T15:04:05", full, location)
	if err != nil {
		zap.L().Error("Couldn't parse time", zap.Error(err))
	}

	return t.Format(time.RFC1123)
}
