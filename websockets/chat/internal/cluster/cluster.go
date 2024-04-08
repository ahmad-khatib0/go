package cluster

import (
	"encoding/gob"
	"errors"
	"sort"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/concurrency"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
)

// NewCluster() returns snowflake worker id (pass nil if you don't want to use cluster)
// Cluster won't be started here yet.
func NewCluster(ca ClusterArgs) (models.Cluster, int, error) {
	if ca.Cfg == nil {
		return nil, 1, nil
	}

	if ca.Cfg.MainName == "" {
		// c.Logger.Info("Cluster: running as a standalone server.")
		return nil, 1, nil
	}

	// INFO: gob is like json, xml, protobuf, but But for a Go-specific environment, such as
	// communicating between two servers written in Go,
	gob.Register([]any{})
	gob.Register(map[string]string{})
	gob.Register(map[string]int{})
	gob.Register(models.MsgAccessMode{})

	res := Cluster{
		Logger:          ca.Logger,
		thisNodeName:    ca.Cfg.MainName,
		fingerprint:     time.Now().Unix(),
		nodes:           make(map[string]*ClusterNode),
		proxyEventQueue: concurrency.NewGoRoutinePool(len(ca.Cfg.Nodes) * 5),
	}

	var nodeNames []string
	for _, host := range ca.Cfg.Nodes {
		nodeNames = append(nodeNames, host.Name)
		if host.Name == ca.Cfg.MainName {
			// Don't create a cluster member for this local instance
			res.listenOn = host.Addr
			continue
		}

		// INFO: An example of session multiplexing—a single computer with one
		// IP address has several websites open at once

		res.nodes[host.Name] = &ClusterNode{
			address: host.Addr,
			name:    host.Name,
			done:    make(chan bool, 1),
			msess:   make(map[string]struct{}),
		}
	}

	if len(res.nodes) == 0 {
		return nil, 1, errors.New("invalid cluster size: Cluster needs at least two nodes.")
	}

	// TODO: add the failoverInit here

	sort.Strings(nodeNames)
	wid := sort.SearchStrings(nodeNames, ca.Cfg.MainName) + 1
	ca.Stats.IntStatsSet(constants.StatsClusterTotalNodes, int64(len(res.nodes)+1))

	return res, wid, nil
}
