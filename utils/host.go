package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	// HostIPEnvVar is the environment variable to override host's chosen IP address.
	HostIPEnvVar = "EGOIST_HOST_IP"
)

func GetHostAddress() (string, error) {
	if val, ok := os.LookupEnv(HostIPEnvVar); ok && val != "" {
		return val, nil
	}

	// Use udp so no handshake is made.
	// Any IP can be used, since connection is not established, but we used a known DNS IP.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return "", fmt.Errorf("error getting interface IP addresses:%w", err)
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}

		return "", errors.New("could not determine host IP address")
	}

	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
