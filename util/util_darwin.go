// build +darwin

package util

import (
	"syscall"
	"unsafe"
)

func SysctlUint64(name string) (uint64, error) {
	v, err := syscall.Sysctl(name)
	if err != nil {
		return 0, err
	}
	buf := []byte(v)
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}
