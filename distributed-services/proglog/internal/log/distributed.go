package log

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"

	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type DistributedLog struct {
	config Config
	log    *Log
	raft   *raft.Raft
}

func NewDistributedLog(dataDir string, config Config) (*DistributedLog, error) {
	l := &DistributedLog{config: config}

	if err := l.setupLog(dataDir); err != nil {
		return nil, err
	}

	if err := l.setupRaft(dataDir); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *DistributedLog) setupLog(dataDir string) error {
	logDir := filepath.Join(dataDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	var err error
	l.log, err = NewLog(logDir, l.config)
	return err
}

// setupRaft(dataDir string) configures and creates the server’s Raft instance.
func (l *DistributedLog) setupRaft(dataDir string) error {
	fsm := &fsm{log: l.log} // creating finite-state-machine (FSM)

	logDir := filepath.Join(dataDir, "raft", "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logConfig := l.config
	logConfig.Segment.InitialOffset = 1 //  initial offset to 1, as required by Raft

	logStore, err := newLogStore(logDir, logConfig)
	if err != nil {
		return err
	}

	// 	The stable store is a key-value store where Raft stores important metadata,
	// like the server’s current term or the candidate the server voted for.
	// Bolt is an embedded and persisted key-value database for Go we’ve used as our stable store.
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft", "stable"))
	if err != nil {
		return err
	}
	retain := 1 //  we’ll keep one snapshot

	// Raft snapshots to recover and restore data efficiently, when necessary, like if your server’s EC2 instance
	// failed and an autoscaling group brought up another instance for the Raft server.
	// Rather than streaming all the data from the Raft leader, the new server would restore from the snapshot
	// and then get the latest changes from the leader. This is more efficient and less taxing on the leader
	// YOU wanna to snapshot frequently to minimize the difference between the data in snapshots and on the leader
	snapshotStore, err := raft.NewFileSnapshotStore(filepath.Join(dataDir, "raft"), retain, os.Stderr)

	if err != nil {
		return err
	}
	maxPool := 5
	timeout := 10 * time.Second
	transport := raft.NewNetworkTransport(l.config.Raft.StreamLayer, maxPool, timeout, os.Stderr)

	config := raft.DefaultConfig()
	config.LocalID = l.config.Raft.LocalID
	// LocalID is the unique ID for this server and it’s the only config field we must set; rest are optional,

	if l.config.Raft.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = l.config.Raft.HeartbeatTimeout
	}

	if l.config.Raft.ElectionTimeout != 0 {
		config.ElectionTimeout = l.config.Raft.ElectionTimeout
	}

	if l.config.Raft.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = l.config.Raft.LeaderLeaseTimeout
	}

	if l.config.Raft.CommitTimeout != 0 {
		config.CommitTimeout = l.config.Raft.CommitTimeout
	}

	// create the Raft instance and bootstrap the cluster:
	l.raft, err = raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return err
	}

	hasState, err := raft.HasExistingState(logStore, stableStore, snapshotStore)
	if err != nil {
		return err
	}

	//  ┌───────────────────────────────────────────────────────────────────────────────┐
	//    Generally you’ll bootstrap a server configured with itself as the only voter,
	//    wait until it becomes the leader, and then tell the leader to add more servers
	//    to the cluster. The subsequently added servers don’t bootstrap.
	//  └───────────────────────────────────────────────────────────────────────────────┘
	if l.config.Raft.Bootstrap && !hasState {
		config := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			}},
		}
		err = l.raft.BootstrapCluster(config).Error()
	}
	return err
}
