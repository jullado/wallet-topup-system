package config

import "time"

func NewAppInitTime() {
	// Set local timezone to Asia/Bangkok (UTC+7)
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}
