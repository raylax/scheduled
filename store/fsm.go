package store

import (
	"bytes"
	"github.com/hashicorp/raft"
	"github.com/raylax/scheduled/codec"
	"io"
)

type CommandRequest struct {
}

type CommandResponse struct {
	Err error
}

var CommandResponseOK CommandResponse

type FSM struct {
	data *Data
}

func NewFSM(data *Data) *FSM {
	return &FSM{data: data}
}

func (f *FSM) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		decoder := codec.NewDecoder(bytes.NewReader(log.Data))
		var command CommandRequest
		err := decoder.Decode(&command)
		if err != nil {
			return err
		}
		return f.applyCommand(command)
	}
	return nil
}

func (f *FSM) applyCommand(command CommandRequest) CommandResponse {
	return CommandResponseOK
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{data: f.data}, nil
}

func (f *FSM) Restore(snapshot io.ReadCloser) error {
	d, err := RestoreData(snapshot)
	if err != nil {
		return err
	}
	f.data = d
	return nil
}
