package hg

import "net"

func GetIPFromDial(dest string) net.IP {
	conn, err := net.Dial("udp", dest)
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
