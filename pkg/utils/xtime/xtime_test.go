package xtime

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	t1 := Time{now}
	fmt.Println(t1)

	fmt.Println(Now())
	t2, _ := Parse(TimeFormat, "2022-04-02 15:00:20")
	fmt.Println(t2)
}
