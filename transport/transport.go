package transport

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"net"
	"time"
)

func NewTransport(listen string, logger hclog.Logger) (*raft.NetworkTransport, error) {
	address, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransportWithLogger(
		address.String(),
		address,
		3,
		10*time.Second,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return transport, nil
}
