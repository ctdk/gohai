// +build linux darwin

package filesystem

import (
	"os/exec"
	"regexp"
	"strings"
)

func getFileSystemInfo() (interface{}, error) {

	/* Grab filesystem data from df	*/
	out, err := exec.Command("df", dfOptions...).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	fileSystemInfo := make(map[string]interface{}, len(lines)-2)
	re := regexp.MustCompile(dfRegexp)
	for _, line := range lines[1:] {
		values := re.FindStringSubmatch(line)
		if len(values) == expectedLength {
			fileSystemInfo[values[1]] = updatefileSystemInfo(values)
		}
	}
	out, err = exec.Command("mount").Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(string(out), "\n")
	mre := regexp.MustCompile(mountRegexp)
	for _, line := range lines {
		values := mre.FindStringSubmatch(line)
		if len(values) == mountLength {
			fileSystemInfo[values[1]] = setMountInfo(values, fileSystemInfo[values[1]].(map[string]interface{}))
		}
	}

	return fileSystemInfo, nil
}
