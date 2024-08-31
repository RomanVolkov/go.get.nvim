package utils_test

import (
	"fmt"
	"main/utils"
	"testing"
	"time"
)

func TestMakeTimeRange(t *testing.T) {
	start, err := time.Parse(time.RFC3339Nano, "2024-01-04T00:00:00.000000Z")
	if err != nil {
		t.Error(err)
		return
	}

	end := start.Add(time.Hour * 10)

	times := utils.MakeTimeRange(start, end, time.Hour*2)
	fmt.Println(len(times))
	if len(times) != 5 {
		t.Error("incorrect slising")
	}

	diff := times[1].Sub(start)
	if diff != time.Hour*2 {
		t.Error("incorrect slicing duration")
	}
}
