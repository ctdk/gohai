package filesystem

import (
	"strconv"
	"strings"
)

const expectedLength = 10
const mountLength = 5

var dfOptions = []string{"-k"}
var dfRegexp = `^(.+?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+\%)\s+(\d+)\s+(\d+)\s+(\d+\%)\s+(.+)$`
var mountRegexp = `^(.+?) on (.+?) \((.+?), (.+?)\)$`

func updatefileSystemInfo(values []string) map[string]interface{} {
	iused, _ := strconv.Atoi(values[6])
	iavail, _ := strconv.Atoi(values[7])
	totalInodes := iused + iavail
	return map[string]interface{}{
		"block_size": "1024",
		"kb_size":    values[2],
		"kb_used":    values[3],
		"kb_available": values[4],
		"percent_used": values[5],
		"inodes_used": values[6],
		"inodes_available": values[7],
		"inodes_percent_used": values[8],
		"total_inodes": strconv.Itoa(totalInodes),
		"mount": values[9],
	}
}

func setMountInfo(values []string, entry map[string]interface{}) map[string]interface{} {
	if entry == nil {
		entry = make(map[string]interface{})
		entry["mount"] = values[2]
	}
	entry["fs_type"] = values[3]
	entry["mount_options"] = strings.Split(values[4], ", ")
	return entry
}
