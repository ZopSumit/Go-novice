package service

import (
	"strings"
	"time"
)

func DateSubstract(d1 string) int {
	d1Slice := strings.Split(d1, "/")

	newDate := d1Slice[2] + "-" + d1Slice[1] + "-" + d1Slice[0]
	myDate, err := time.Parse("2006-01-02", newDate)

	if err != nil {
		//panic(err)
		return 0
	}

	return int(time.Now().Unix() - myDate.Unix())
}
