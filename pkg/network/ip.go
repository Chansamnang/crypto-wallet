package network

import (
	"fmt"
	"net"

	"wallet/pkg/constant"
)

func GetRpcRegisterIP(configIP string) (string, error) {
	registerIP := configIP

	if registerIP == "" {
		ip, err := GetLocalIP()
		if err != nil {
			return "", err
		}
		registerIP = ip
	}
	return registerIP, nil
}

func GetListenIP(configIP string) string {
	if configIP == "" {
		return constant.LocalHost
	} else {
		return configIP
	}
}

var ServerIP = ""

func GetLocalIP() (string, error) {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				if ipnet.IP.IsLoopback() {
					continue
				}
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	if len(ips) > 0 {
		return ips[0], nil
	}

	return "", fmt.Errorf("no ip found")
}
