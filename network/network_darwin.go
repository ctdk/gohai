// +build darwin

package network

import (
	"bufio"
	"bytes"
	"net"
	"os/exec"
	"regexp"
)

func getNetworkCounters() (map[string]interface{}, error) {
	n, err := exec.Command("netstat", "-i", "-d", "-l", "-b", "-n").Output()
	if err != nil {
		return nil, err
	}
	ifaces := make(map[string]interface{})
	nread := bufio.NewScanner(bytes.NewBuffer(n))
	re1 := regexp.MustCompile(`^([a-zA-Z0-9\.\:\-\*]+)\s+\d+\s+\<[a-zA-Z0-9\#]+\>\s+([a-f0-9\:]+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	re2 := regexp.MustCompile(`^([a-zA-Z0-9\.\:\-\*]+)\s+\d+\s+\<[a-zA-Z0-9\#]+\>(\s+)(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
	for nread.Scan() {
		l := nread.Text()
		var matchLine []string
		if matchLine = re1.FindStringSubmatch(l); matchLine == nil {
			matchLine = re2.FindStringSubmatch(l)
		}
		if matchLine != nil {
			if i, _ := net.InterfaceByName(matchLine[1]); i == nil {
				continue
			}
			rx := map[string]interface{}{"bytes": matchLine[5], "packets": matchLine[3], "errors": matchLine[4], "drop": 0, "frame": 0, "compressed": 0, "multicast": 0}
			tx := map[string]interface{}{"bytes": matchLine[8], "packets": matchLine[6], "errors": matchLine[7], "overrun": 0, "collisions": matchLine[9], "carrier": 0, "compressed": 0}
			ifc := make(map[string]interface{})
			ifc["rx"] = rx
			ifc["tx"] = tx
			ifaces[matchLine[1]] = ifc
		}
	}
	f := map[string]interface{}{"interfaces": ifaces}
	fullInfo := map[string]interface{}{"network": f}
	return fullInfo, nil
}
