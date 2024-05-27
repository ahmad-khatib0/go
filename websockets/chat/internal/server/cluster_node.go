package server

import (
	"errors"
	"net"
	"net/rpc"
	"time"
)

func (n *ClusterNode) asyncRpcLoop() {
	for call := range n.rpcDone {
		n.handleRpcResponse(call)
	}
}

func (n *ClusterNode) p2mSenderLoop() {
	for req := range n.p2mSender {
		if req == nil {
			// Stop
			return
		}

		if err := n.proxyToMaster(req); err != nil {
			globals.l.Sugar().Warnf("p2mSenderLoop: call failed", n.name, err)
		}
	}
}

func (n *ClusterNode) call(proc string, req, resp any) error {
	if !n.connected {
		return errors.New("cluster: node '" + n.name + "' not connected")
	}

	if err := n.endpoint.Call(proc, req, resp); err != nil {
		globals.l.Sugar().Warnf("cluster: call failed", n.name, err)

		n.lock.Lock()
		if n.connected {
			n.endpoint.Close()
			n.connected = false
			globals.stats.IntStatsInc("LiveClusterNodes", -1)
			go n.reconnect()
		}
		n.lock.Unlock()
		return err
	}

	return nil
}

func (n *ClusterNode) handleRpcResponse(call *rpc.Call) {
	if call.Error != nil {
		globals.l.Sugar().Warnf("cluster: %s call failed: %s", call.ServiceMethod, call.Error)
		n.lock.Lock()

		if n.connected {
			n.endpoint.Close()
			n.connected = false

			globals.stats.IntStatsInc("LiveClusterNodes", -1)
			go n.reconnect()
		}

		n.lock.Unlock()
	}
}

func (n *ClusterNode) callAsync(proc string, req, resp any, done chan *rpc.Call) *rpc.Call {
	if done != nil && cap(done) == 0 {
		globals.l.Panic("cluster: RPC done channel is unbuffered")
	}

	if !n.connected {
		call := &rpc.Call{
			ServiceMethod: proc,
			Args:          req,
			Reply:         resp,
			Error:         errors.New("cluster: node '" + n.name + "' not connected"),
			Done:          done,
		}
		if done != nil {
			done <- call
		}
		return call
	}

	var responseChan chan *rpc.Call
	if done != nil {
		// Make a separate response callback if we need to notify the caller.
		myDone := make(chan *rpc.Call, 1)
		go func() {
			call := <-myDone
			n.handleRpcResponse(call)
			if done != nil {
				done <- call
			}
		}()
		responseChan = myDone
	} else {
		responseChan = n.rpcDone
	}

	call := n.endpoint.Go(proc, req, resp, responseChan)

	return call
}

// proxyToMaster forwards request from topic proxy to topic master.
func (n *ClusterNode) proxyToMaster(msg *ClusterReq) error {
	msg.Node = globals.cluster.thisNodeName
	var rejected bool

	err := n.call("Cluster.TopicMaster", msg, &rejected)
	if err == nil && rejected {
		err = errors.New("cluster: topic master node out of sync")
	}
	return err
}

// proxyToMaster forwards request from topic proxy to topic master.
func (n *ClusterNode) proxyToMasterAsync(msg *ClusterReq) error {
	select {
	case n.p2mSender <- msg:
		return nil
	default:
		return errors.New("cluster: load exceeded")
	}
}

// masterToProxyAsync forwards response from topic master to topic proxy
// in a fire-and-forget manner.
func (n *ClusterNode) masterToProxyAsync(msg *ClusterResp) error {
	var unused bool
	if c := n.callAsync("Cluster.TopicProxy", msg, &unused, nil); c.Error != nil {
		return c.Error
	}
	return nil
}

// route routes server message within the cluster.
func (n *ClusterNode) route(msg *ClusterRoute) error {
	var unused bool
	return n.call("Cluster.Route", msg, &unused)
}

// Handle outbound node communication: read messages from the channel, forward to remote nodes.
//
// FIXME(gene): this will drain the outbound queue in case of a failure: all unprocessed

// messages will be dropped. Maybe it's a good thing, maybe not.
func (n *ClusterNode) reconnect() {
	var reconnTicker *time.Ticker

	// Avoid parallel reconnection threads.
	n.lock.Lock()
	if n.reconnecting {
		n.lock.Unlock()
		return
	}
	n.reconnecting = true
	n.lock.Unlock()

	count := 0
	for {
		// Attempt to reconnect right away
		if conn, err := net.DialTimeout("tcp", n.address, clusterNetworkTimeout); err == nil {
			if reconnTicker != nil {
				reconnTicker.Stop()
			}

			n.lock.Lock()
			n.endpoint = rpc.NewClient(conn)
			n.connected = true
			n.reconnecting = false
			n.lock.Unlock()
			globals.stats.IntStatsInc("LiveClusterNodes", 1)
			globals.l.Sugar().Infof("cluster: connected to", n.name)

			// Send this node credentials to the new node.
			var unused bool
			n.call(
				"Cluster.Ping",
				&ClusterPing{Node: globals.cluster.thisNodeName, Fingerprint: globals.cluster.fingerprint},
				&unused,
			)
			return
		} else if count == 0 {
			reconnTicker = time.NewTicker(clusterDefaultReconnectTime)
		}

		count++

		select {
		case <-reconnTicker.C:
			// Wait for timer to try to reconnect again. Do nothing if the timer is inactive.
		case <-n.done:
			// Shutting down
			globals.l.Sugar().Infof("cluster: shutdown started at node", n.name)
			reconnTicker.Stop()
			if n.endpoint != nil {
				n.endpoint.Close()
			}
			n.lock.Lock()
			n.connected = false
			n.reconnecting = false
			n.lock.Unlock()
			globals.l.Sugar().Infof("cluster: shut down completed at node", n.name)
			return
		}
	}
}

func (n *ClusterNode) stopMultiplexingSession(msess *Session) {
	if msess == nil {
		return
	}

	msess.stopSession(nil)
	n.lock.Lock()
	delete(n.msess, msess.sid)
	n.lock.Unlock()
}
