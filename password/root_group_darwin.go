// build +darwin

package password

import (
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
)

func lookupRootGroup() (string, error) {
	rootUser, err := user.Lookup("root")
	if err != nil {
		return "", err
	}
	groupInfo, err := exec.Command("dscacheutil", "-q", "group", "-a", "gid", rootUser.Gid).Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^name: (\w+)`)
	s := re.FindStringSubmatch(string(groupInfo))
	if len(s) < 2 {
		err = fmt.Errorf("could not find group name from dscacheutil")
		return "", err
	}
	return s[1], nil
}
