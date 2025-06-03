package util

import (
	"strings"
	"time"
)

func FormatIndonesianDate(date string) (string, error) {
	monthMap := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	dateObj, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}

	indonesianDate := dateObj.Format("2 January 2006")
	for engMonth, indMonth := range monthMap {
		indonesianDate = strings.Replace(indonesianDate, engMonth, indMonth, 1)
	}

	return indonesianDate, nil
}
