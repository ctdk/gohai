// +build linux darwin

package network

func getNetworkInfo() (networkInfo map[string]interface{}, err error) {
	networkInfo = make(map[string]interface{})

	ifaces, err := networkInterfaces()
	if err != nil {
		return nil, err
	}
	networkInfo["interfaces"] = ifaces

	settings, err := settings()
	if err != nil {
		return nil, err
	}
	networkInfo["settings"] = settings

	return
}

func getTopLevel() (map[string]interface{}, error) {
	networkInfo := make(map[string]interface{})
	macaddress, err := macAddress()
	if err != nil {
		return networkInfo, err
	}
	networkInfo["macaddress"] = macaddress

	ipAddress, err := externalIpAddress()
	if err != nil {
		return networkInfo, err
	}
	networkInfo["ipaddress"] = ipAddress

	ipAddressV6, err := externalIpv6Address()
	if err != nil {
		return networkInfo, err
	}
	networkInfo["ipaddressv6"] = ipAddressV6
	return networkInfo, nil
}

func getNetworkCounters() (map[string]interface{}, error) {
	ifaces := make(map[string]interface{})
}
