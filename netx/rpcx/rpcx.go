package rpcx

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

// DialTimeout acts like Dial but takes a timeout. Returns a new rpc.Client to
// handle requests to the set of services at the other end of the connection.
func DialTimeout(network, address string, timeout time.Duration) (*rpc.Client, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), err
}
