// +build darwin

package cpu

import (
	"strings"
	"syscall"
	"unsafe"
)

func getCpuInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	cpuKeys := map[string]string{"vendor": "vendor_id", "brand_string": "model_name"}
	for sysctlName, gohaiName := range cpuKeys {
		k, err := syscall.Sysctl("machdep.cpu." + sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	cpuKeyInts := map[string]string{"model": "model", "family": "family", "stepping": "stepping"}
	for sysctlName, gohaiName := range cpuKeyInts {
		k, err := syscall.SysctlUint32("machdep.cpu." + sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	hwKeyInts := map[string]string{"physicalcpu": "real", "logicalcpu": "total", "cpufrequency": "mhz"}
	for sysctlName, gohaiName := range hwKeyInts {
		k, err := sysctluint64("hw." + sysctlName)
		if err != nil {
			return nil, err
		}
		info[gohaiName] = k
	}
	info["mhz"] = info["mhz"].(uint64) / 1000000

	cpuFlags, err := syscall.Sysctl("machdep.cpu.features")
	if err != nil {
		return nil, err
	}
	info["flags"] = strings.Split(strings.ToLower(cpuFlags), " ")

	return info, nil
}

func sysctluint64(name string) (uint64, error) {
	v, err := syscall.Sysctl(name)
	if err != nil {
		return 0, err
	}
	buf := []byte(v)
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}
