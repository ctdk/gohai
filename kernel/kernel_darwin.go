// +build darwin

package kernel

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

func getKernelInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})
	// still need modules
	kernKeys := map[string]string{"ostype": "name", "osrelease": "release", "version": "version"}
	for sysctlName, gohaiName := range kernKeys {
		k, err := syscall.Sysctl("kern." + sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	k, err := syscall.Sysctl("hw.machine")
	if err != nil {
		return nil, err
	}
	info["machine"] = k
	info["os"] = info["name"]
	modinfo, err := exec.Command("kextstat", "-k", "-l").Output()
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(\d+)\s+(\d+)\s+0x[0-9a-f]+\s+0x([0-9a-f]+)\s+0x[0-9a-f]+\s+([a-zA-Z0-9\.]+) \(([0-9\.]+)\)`)
	modbuf := bytes.NewBuffer(modinfo)
	modread := bufio.NewScanner(modbuf)
	info["modules"] = make(map[string]map[string]interface{})
	for modread.Scan() {
		s := re.FindStringSubmatch(modread.Text())
		if len(s) < 6 {
			continue
		}
		size, err := strconv.ParseInt(s[3], 16, 64)
		if err != nil {
			return nil, err
		}
		info["modules"].(map[string]map[string]interface{})[s[4]] = map[string]interface{}{"version": s[5], "size": size, "index": s[1], "refcount": s[2]}
	}

	fullInfo := map[string]interface{}{"kernel": info}
	return fullInfo, nil
}
