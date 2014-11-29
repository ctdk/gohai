// build !windows

package password

import (
	"bufio"
	"os"
	"strings"
)

func getPasswordInfo() (map[string]interface{}, error) {
	var err error
	info := make(map[string]interface{})
	info["passwd"], err = passwdInfo()
	if err != nil {
		return nil, err
	}
	info["group"], err = groupInfo()
	if err != nil {
		return nil, err
	}
	fullInfo := map[string]interface{}{"etc": info}
	return fullInfo, nil
}

func passwdInfo() (map[string]interface{}, error) {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	pinfo := make(map[string]interface{})
	passwds := bufio.NewScanner(f)
	for passwds.Scan() {
		p := passwds.Text()
		if p[0] == '#' {
			continue
		}
		fields := strings.Split(p, ":")
		if len(fields) != 7 {
			continue
		}
		l := make(map[string]interface{})
		l["dir"] = fields[5]
		l["uid"] = fields[2]
		l["gid"] = fields[3]
		l["shell"] = fields[6]
		l["gecos"] = fields[4]
		pinfo[fields[0]] = l
	}
	return pinfo, nil
}

func groupInfo() (map[string]interface{}, error) {
	f, err := os.Open("/etc/group")
	if err != nil {
		return nil, err
	}
	ginfo := make(map[string]interface{})
	groups := bufio.NewScanner(f)
	for groups.Scan() {
		g := groups.Text()
		if g[0] == '#' {
			continue
		}
		fields := strings.Split(g, ":")
		if len(fields) != 4 {
			continue
		}
		l := make(map[string]interface{})
		l["gid"] = fields[2]
		var members []string
		if fields[3] != "" {
			members = strings.Split(fields[3], ",")
		} else {
			members = make([]string, 0)
		}
		l["members"] = members
		ginfo[fields[0]] = l
	}
	return ginfo, nil
}
