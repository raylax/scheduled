package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/raylax/scheduled/node"
	"github.com/raylax/scheduled/store"
	"github.com/raylax/scheduled/transport"
	"github.com/raylax/scheduled/types"
	"github.com/raylax/scheduled/util"
	"os"
	"time"
)

const listen = "127.0.0.1:7001"
const dataPath = "data"

func main() {

	if !util.IsExists(dataPath) {
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "RAFT",
		Level: hclog.LevelFromString("debug"),
	})

	data := store.NewData()
	fsm := store.NewFSM(data, logger)

	t, err := transport.NewTransport(listen, logger)
	if err != nil {
		panic(err)
	}

	id := raft.ServerID(listen)

	opts := node.Options{
		ID:                id,
		Listen:            listen,
		DataPath:          dataPath,
		Logger:            logger,
		Transport:         t,
		SnapshotInterval:  10 * time.Second,
		SnapshotThreshold: 2,
	}
	n, err := node.New(opts, fsm)
	if err != nil {
		panic(err)
	}

	config := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      id,
				Address: t.LocalAddr(),
			},
		},
	}

	cluster := n.Raft.BootstrapCluster(config)


	if err := cluster.Error(); err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			logger.Info("store", "data", data)
		}
	}()

	for {
		leader := <-n.LeaderCh
		logger.Info("state change", "id", id, "leader", leader)
		if leader {
			go func() {
				bytes, err := types.EncodeCommandRequest(types.CommandTypeSet, &types.CommandSet{
					Key:   "a",
					Value: "1",
				})
				if err != nil {
					logger.Warn("encode command error", "err", err)
					return
				}

				f := n.Raft.Apply(bytes, 0)
				logger.Info("apply", "response", f.Response(), "err", f.Error())
			}()
		}
	}

}
