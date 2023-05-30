package main

import (
	"testing"
	"time"
)

func Test(t *testing.T) {

	startTime := 1685081469000
	format := "2006-01-02 15:04:05"

	t1 := time.UnixMilli(int64(startTime)).Format(format)

	t.Log(t1)
}
