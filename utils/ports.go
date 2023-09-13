package utils

import (
	"crypto/sha256"
	"net"
	"os"
	"strconv"
	"strings"
)

func GetStablePort(start int, appId string) (int, error) {
	parts := make([]string, 4)
	parts[0] = appId
	parts[1], _ = os.Getwd()
	parts[2] = strconv.Itoa(os.Getuid())
	parts[3], _ = os.Hostname()

	base := []byte(strings.Join(parts, "|"))

	h := sha256.Sum256(base)

	//取前11位(0-2047)
	rnd := int(h[0]) + int(h[1]>>5)<<8
	port := start + rnd

	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err == nil {
		_ = l.Close()
		return port, nil
	}

	return GetFreePort()
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
