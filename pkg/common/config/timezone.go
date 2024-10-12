package config

import (
	"fmt"
	"time"
)

const (
	DefaultTimezone        = "Asia/Manila"
	DefaultTimeOffsetHours = 8 // +8:00,
)

var (
	timezone        *time.Location
	timeOffsetHours int
)

func initTimezone() {
	var err error
	timezone, err = time.LoadLocation(Config.App.Timezone)
	if err != nil {
		fmt.Printf("invalid timezone from config: %q, %s", Config.App.Timezone, err.Error())

		// try use default timezone,
		timezone, err = time.LoadLocation(DefaultTimezone)
		if err != nil {
			fmt.Printf("invalid default timezone: %q, %s", DefaultTimezone, err.Error())
			return
		}
		timeOffsetHours = DefaultTimeOffsetHours
	}

	_, timeOffsetSec := time.Now().In(timezone).Zone()
	// TODO: 目前只支持 整数 offset.
	timeOffsetHours = timeOffsetSec / int(time.Hour.Seconds())

	fmt.Printf("[init timezone] timezone = %q, timeOffsetHours = %d\n", timezone.String(), timeOffsetHours)
}

func GetTimeOffsetHours() int {
	return timeOffsetHours
}

func GetTimezoneName() string {
	return timezone.String()
}

func GetTimezoneLoc() *time.Location {
	return timezone
}
