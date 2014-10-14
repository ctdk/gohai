// +build darwin

package platform

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"time"
)

func getArchInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})
	// still need modules
	kernKeys := map[string]string{"ostype": "os", "osrelease": "os_version"}
	for sysctlName, gohaiName := range kernKeys {
		k, err := syscall.Sysctl("kern." + sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	info["platform"] = "mac_os_x"
	info["platform_family"] = "mac_os_x"
	// This may prove worth generalizing
	tval := new(syscall.Timeval)
	u, err := syscall.Sysctl("kern.boottime")
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte(u))
	binary.Read(buf, binary.LittleEndian, tval)
	info["uptime_seconds"] = int64(time.Since(time.Unix(tval.Unix())).Seconds())

	info["ohai_time"] = time.Now().Unix()

	return info, nil
}
