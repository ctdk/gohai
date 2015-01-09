// +build darwin

package platform

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-chef/gohai/util"
	"os/exec"
	"regexp"
	"strings"
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
	info["platform_family"] = "mac_os_x"
	p, err := exec.Command("/usr/bin/sw_vers").Output()
	if err != nil {
		return nil, err
	}
	pread := bufio.NewScanner(bytes.NewBuffer(p))
	for pread.Scan() {
		l := pread.Text()
		switch {
		case strings.HasPrefix(l, "ProductName:"):
			re := regexp.MustCompile(`ProductName:\s+(.+)$`)
			ml := re.FindStringSubmatch(l)
			m := ml[1]
			m = strings.ToLower(m)
			info["platform"] = strings.Replace(m, " ", "_", -1)
		case strings.HasPrefix(l, "ProductVersion:"):
			re := regexp.MustCompile(`ProductVersion:\s+(.+)$`)
			ml := re.FindStringSubmatch(l)
			info["platform_version"] = ml[1]

		case strings.HasPrefix(l, "BuildVersion:"):
			re := regexp.MustCompile(`BuildVersion:\s+(.+)$`)
			ml := re.FindStringSubmatch(l)
			info["platform_build"] = ml[1]
		}
	}
	// This may prove worth generalizing
	tval := new(syscall.Timeval)
	u, err := syscall.Sysctl("kern.boottime")
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte(u))
	binary.Read(buf, binary.LittleEndian, tval)
	ut := time.Since(time.Unix(tval.Unix()))
	info["uptime_seconds"] = int64(ut.Seconds())
	info["uptime"] = util.DurationToHuman(ut)

	info["ohai_time"] = fmt.Sprintf("%f", float64(time.Now().UnixNano())/float64(time.Second))

	return info, nil
}
