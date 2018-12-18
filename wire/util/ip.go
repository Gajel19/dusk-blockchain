package util

import (
	"net"
)

// GetLocalIP will return the machine's external IP address.
// https://stackoverflow.com/a/37382208
func GetLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
