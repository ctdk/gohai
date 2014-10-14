package password

import (

)

type Password struct{}
const name = "etc"

func (p *Password) Name() string {
	return name
}

func (p *Password) Collect() (interface{}, error) {
	result, err := getPasswordInfo()
	return result, err
}

type TopLevel struct{}

func (t *TopLevel) Name() string {
	return "top_level"
}

func (t *TopLevel) Collect() (interface{}, error) {
	result, err := getTopLevel()
	return result, err
}
