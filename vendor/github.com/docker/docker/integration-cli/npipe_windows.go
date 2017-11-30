package main

import (
	"net"
	"time"
)

func npipeDial(path string, timeout time.Duration) (net.Conn, error) {
	return winio.DialPipe(path, &timeout)
}
