package store

import (
	"github.com/hashicorp/raft"
)

type Snapshot struct {
	data *Data
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	return s.data.Persist(sink)
}

func (s *Snapshot) Release() {
}
