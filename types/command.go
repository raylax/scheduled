package types

import (
	"bytes"
	"github.com/raylax/scheduled/codec"
)

type CommandType int

const (
	CommandTypeNone CommandType = iota
	CommandTypeSet  CommandType = iota
)

type CommandRequest struct {
	Type CommandType
	Data []byte
}

type CommandSet struct {
	Key   string
	Value string
}

type CommandResponse struct {
	Err error
}

var CommandResponseOK CommandResponse

func EncodeCommandRequest[T any](t CommandType, command T) ([]byte, error) {
	buf := bytes.Buffer{}
	err := codec.NewEncoder(&buf).Encode(&command)
	if err != nil {
		return nil, err
	}

	var req = &CommandRequest{
		Type: t,
		Data: buf.Bytes(),
	}

	buf = bytes.Buffer{}
	err = codec.NewEncoder(&buf).Encode(req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeCommandRequest(b []byte) (*CommandRequest, error) {
	decoder := codec.NewDecoder(bytes.NewReader(b))
	var command CommandRequest
	err := decoder.Decode(&command)
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func DecodeCommand[T any](b []byte) (*T, error) {
	var v T
	err := codec.NewDecoder(bytes.NewReader(b)).Decode(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
