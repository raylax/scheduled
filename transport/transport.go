package transport

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"net"
	"time"
)

func NewTransport(address string, logger hclog.Logger) (*raft.NetworkTransport, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransportWithLogger(
		addr.String(),
		addr,
		3,
		10*time.Second,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return transport, nil
}
