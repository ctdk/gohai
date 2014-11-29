// +build linux

package platform

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var unameOptions = []string{"-s", "-n", "-r", "-m", "-p", "-i", "-o"}

func getArchInfo() (archInfo map[string]interface{}, err error) {
	archInfo = make(map[string]interface{})

	out, err := exec.Command("uname", unameOptions...).Output()
	if err != nil {
		return nil, err
	}
	line := fmt.Sprintf("%s", out)
	values := regexp.MustCompile(" +").Split(line, 7)
	updateArchInfo(archInfo, values)

	out, err = exec.Command("uname", "-v").Output()
	if err != nil {
		return nil, err
	}
	archInfo["kernel_version"] = strings.Trim(string(out), "\n")

	return
}

func updateArchInfo(archInfo map[string]interface{}, values []string) {
	archInfo["kernel_name"] = values[0]
	archInfo["kernel_release"] = values[2]
	archInfo["machine"] = values[3]
	archInfo["processor"] = values[4]
	archInfo["hardware_platform"] = values[5]
	archInfo["os"] = strings.Trim(values[6], "\n")
}
