package helpers

import "net"

func GetServerIp() (ip string) {
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return
}
