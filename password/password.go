package password

import ()

type Password struct{}

const name = "etc"

func (p *Password) Name() string {
	return name
}

func (p *Password) Collect() (interface{}, error) {
	result, err := getPasswordInfo()
	if err != nil {
		return nil, err
	}
	topRes, err := getTopLevel()
	if err != nil {
		return nil, err
	}
	for k, v := range topRes {
		result[k] = v
	}
	return result, err
}

func (p *Password) Provides() []string {
	return []string{"etc", "current_user", "root_group", "etc/group", "etc/passwd"}
}
