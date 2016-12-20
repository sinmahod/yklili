package dateutil

import "time"

func GetYMDPathString() string {
	return time.Now().Format("2006/01/02/")
}
