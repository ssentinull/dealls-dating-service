package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDate(date string) time.Time {
	res, _ := time.Parse("2006-01-02", date)

	return res
}

func ParseDateTime(date string) time.Time {
	res, err := time.Parse(time.RFC3339, date)
	if err != nil {
		res, _ = time.Parse("2006-01-02T15:04:05.999999", date)
	}

	return res
}

func GetBeginningOfDayForLocationInUTC(date time.Time, loc *time.Location) time.Time {
	date = date.In(loc)
	dateTruncated := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc).In(loc)

	return dateTruncated.UTC()
}

func GetEndOfDayForLocationInUTC(date time.Time, loc *time.Location) time.Time {
	date = date.In(loc)
	dateTruncated := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, loc).In(loc)

	return dateTruncated.UTC()
}

func GetTableName(query string, time time.Time) string {
	year := strconv.Itoa(time.Year())
	month := fmt.Sprintf("%02d", int(time.Month()))
	day := fmt.Sprintf("%02d", time.Day())
	strReplacer := strings.NewReplacer(
		"YYYYMMDD", year+month+day,
		"YYYYMM", year+month,
		"YYYY", year,
	)
	return strReplacer.Replace(query)
}
