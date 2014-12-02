package network

type Network struct{}

const name = "network"

func (self *Network) Name() string {
	return name
}

func (self *Network) Collect() (result interface{}, err error) {
	result, err = getNetworkInfo()
	return
}

func (self *Network) Provides() []string {
	return []string{"network", "counters", "macaddress", "ipaddress", "ipaddressv6", "network/interfaces", "network/settings"}
}
