package util

import (
	"fmt"
	"time"
)

// MergeMap merges two map[string]interface{} structures together. It's useful
// for adding collected data to the overall map of data to be converted to JSON.
func MergeMap(dst, src map[string]interface{}) error {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			// easiest situation; the destination doesn't have this
			// key.
			dst[k] = v
			continue
		} else {
			// otherwise it's more complicated
			err := merginate(dst[k], v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func merginate(dst, src interface{}) error {
	switch dst := dst.(type) {
	case map[string]interface{}:
		// but what is src?
		switch src := src.(type) {
		case map[string]interface{}:
			err := MergeMap(dst, src)
			if err != nil {
				return err
			}
		default:
			// others later though?
		}
	case nil:
		dst = src
	}
	return nil
}

func DurationToHuman(d time.Duration) string {
	days := d / (time.Hour * 24)
	d = d - (days * time.Hour * 24)
	hours := d / time.Hour
	d = d - hours*time.Hour
	minutes := d / time.Minute
	d = d - minutes*time.Minute
	seconds := d / time.Second

	var dayStr string
	if days == 1 {
		dayStr = "day"
	} else {
		dayStr = "days"
	}
	var uptime string
	switch {
	case days > 0:
		uptime = fmt.Sprintf("%d %s %02d hours %02d minutes %02d seconds", days, dayStr, hours, minutes, seconds)
	case hours > 0:
		uptime = fmt.Sprintf("%02d hours %02d minutes %02d seconds", hours, minutes, seconds)
	case minutes > 0:
		uptime = fmt.Sprintf("%02d minutes %02d seconds", minutes, seconds)
	default:
		uptime = fmt.Sprintf("%02d seconds", seconds)
	}

	return uptime
}
