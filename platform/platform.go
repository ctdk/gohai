package platform

type Platform struct{}

const name = "platform"

func (self *Platform) Name() string {
	return name
}

func (self *Platform) Collect() (result interface{}, err error) {
	result, err = getPlatformInfo()
	return
}

func (self *Platform) Provides() []string {
	return []string{"os", "os_version", "platform", "platform_family", "uptime_seconds", "ohai_time"}
}
