package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/raylax/scheduled/node"
	"github.com/raylax/scheduled/store"
	"github.com/raylax/scheduled/transport"
	"github.com/raylax/scheduled/util"
	"os"
	"time"
)

const address = "127.0.0.1:7001"
const dataPath = "data"

func main() {

	if !util.IsExists(dataPath) {
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	data := store.NewData()
	fsm := store.NewFSM(data)
	logger := hclog.Default()

	t, err := transport.NewTransport(address, logger)
	if err != nil {
		panic(err)
	}

	id := raft.ServerID(address)

	opts := node.Options{
		ID:                id,
		Address:           address,
		DataPath:          dataPath,
		Logger:            logger,
		Transport:         t,
		SnapshotInterval:  20 * time.Second,
		SnapshotThreshold: 2,
	}
	n, err := node.New(opts, fsm)
	if err != nil {
		panic(err)
	}

	//config := raft.Configuration{
	//	Servers: []raft.Server{
	//		{
	//			ID:      id,
	//			Address: t.LocalAddr(),
	//		},
	//	},
	//}
	//
	//cluster := n.Raft.BootstrapCluster(config)
	//if err := cluster.Error(); err != nil {
	//	panic(err)
	//}

	for {
		leader := <-n.Raft.LeaderCh()
		logger.Info("state change", "id", id, "leader", leader)
	}

}
