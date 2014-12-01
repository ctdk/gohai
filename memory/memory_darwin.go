// build +darwin

package memory

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-chef/gohai/util"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func getMemoryInfo() (map[string]interface{}, error) {
	memoryInfo := make(map[string]interface{})

	// can run afoul of 32 bit limits
	total, err := util.SysctlUint64("hw.memsize")
	if err != nil {
		return nil, err
	}
	total = total / 1024 / 1024
	memoryInfo["total"] = fmt.Sprintf("%dMB", total)

	swapout, err := exec.Command("sysctl", "vm.swapusage").Output() //syscall.Sysctl("vm.swapusage")
	if err != nil {
		return nil, err
	}
	swap := regexp.MustCompile("total = ").Split(string(swapout), 2)[1]
	memoryInfo["swap_total"] = strings.Split(swap, " ")[0]
	pagesize, err := syscall.SysctlUint32("hw.pagesize")
	if err != nil {
		return nil, err
	}
	memoryInfo["pagesize"] = pagesize
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		return nil, err
	}
	vmread := bufio.NewScanner(bytes.NewBuffer(out))
	var (
		active         int64
		inactive       int64
		total_consumed int64
	)
	fields := map[string]string{"active": "Pages active:", "inactive": "Pages inactive:", "wired down": "Pages wired down:"}
	re := regexp.MustCompile(`(\d+)\.$`)
	for vmread.Scan() {
		line := vmread.Text()
		for k, v := range fields {
			if strings.HasPrefix(line, v) {
				u := re.FindStringSubmatch(line)
				if len(u) > 1 {
					m, err := strconv.ParseInt(u[1], 10, 64)
					if err != nil {
						return nil, err
					}
					mem := (m * int64(pagesize)) / 1024 / 1024
					total_consumed += mem
					switch k {
					case "active":
						active += mem
					case "inactive":
						inactive += mem
					case "wired down":
						active += mem
					}
				}
			}
		}
	}
	if active > 0 {
		memoryInfo["active"] = fmt.Sprintf("%dMB", active)
	}
	if inactive > 0 {
		memoryInfo["inactive"] = fmt.Sprintf("%dMB", inactive)
	}
	if total_consumed > 0 {
		memoryInfo["free"] = fmt.Sprintf("%dMB", int64(total)-total_consumed)
	}

	fullInfo := map[string]interface{}{"memory": memoryInfo}
	return fullInfo, nil
}
