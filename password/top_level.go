package password

import (
	"os/user"
)

func getTopLevel() (map[string]interface{}, error) {
	curUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	info := make(map[string]interface{})
	info["current_user"] = curUser.Username

	rootGroup, err := lookupRootGroup()
	if err != nil {
		return nil, err
	}
	if rootGroup != "" {
		info["root_group"] = rootGroup
	}
	return info, nil
}
