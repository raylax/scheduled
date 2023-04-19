package store

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/raylax/scheduled/types"
	"io"
)

type FSM struct {
	data   *Data
	logger hclog.Logger
}

func NewFSM(data *Data, logger hclog.Logger) *FSM {
	return &FSM{data: data, logger: logger}
}

func (f *FSM) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		req, err := types.DecodeCommandRequest(log.Data)
		if err != nil {
			return types.CommandResponse{Err: err}
		}
		err = f.applyCommand(req)
		if err != nil {
			return types.CommandResponse{Err: err}
		}
		return types.CommandResponseOK
	}
	return nil
}

func (f *FSM) applyCommand(req *types.CommandRequest) error {
	switch req.Type {
	case types.CommandTypeSet:
		command, err := types.DecodeCommand[types.CommandSet](req.Data)
		if err != nil {
			return err
		}
		return f.applySetCommand(command)
	}
	return nil
}

func (f *FSM) applySetCommand(command *types.CommandSet) error {
	f.data.Set(command.Key, command.Value)
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{data: f.data}, nil
}

func (f *FSM) Restore(snapshot io.ReadCloser) error {
	return f.data.Restore(snapshot)
}
