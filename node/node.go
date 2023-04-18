package node

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"path"
	"time"
)

type Options struct {
	ID        raft.ServerID
	Address   string
	DataPath  string
	Logger    hclog.Logger
	Transport raft.Transport

	SnapshotInterval  time.Duration
	SnapshotThreshold int
}

type Node struct {
	Raft *raft.Raft
}

func New(opts Options, fsm raft.FSM) (*Node, error) {

	logStore, err := boltdb.NewBoltStore(path.Join(opts.DataPath, "log.dat"))
	if err != nil {
		return nil, err
	}

	stableStore, err := boltdb.NewBoltStore(path.Join(opts.DataPath, "stable.dat"))
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStoreWithLogger(opts.DataPath, 1, opts.Logger)
	if err != nil {
		return nil, err
	}

	config := raft.DefaultConfig()
	config.LocalID = opts.ID
	config.Logger = opts.Logger
	config.NotifyCh = make(chan bool)

	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, opts.Transport)

	return &Node{
		Raft: r,
	}, nil
}
