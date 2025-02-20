package netctx

import (
	"net"
)

// Pipe creates piped pair of Conn.
func Pipe() (Conn, Conn) {
	ca, cb := net.Pipe()

	return NewConn(ca), NewConn(cb)
}
