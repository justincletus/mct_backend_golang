package tz

import "os"

func init() {
	os.Setenv("TZ", "Asia/Kolkata")
}
