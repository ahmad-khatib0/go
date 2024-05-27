package server

import (
	"fmt"
	"sync/atomic"
	"time"
)

// queueOut attempts to send a ServerComMessage to a session write loop;
// it fails, if the send buffer is full.
func (s *Session) queueOut(msg *ServerComMessage) bool {
	if s == nil {
		return true
	}

	if atomic.LoadInt32(&s.terminating) > 0 {
		return true
	}

	if s.multi != nil {
		// In case of a cluster we need to pass a copy of the actual session.
		msg.sess = s
		if s.multi.queueOut(msg) {
			s.multi.scheduleClusterWriteLoop()
			return true
		}
		return false
	}

	// Record latency only on {ctrl} messages and end-user sessions.
	if msg.Ctrl != nil && msg.Id != "" {

		if !msg.Ctrl.Timestamp.IsZero() && !s.isCluster() {
			duration := time.Since(msg.Ctrl.Timestamp).Milliseconds()
			globals.stats.HistogramAddSample("RequestLatency", float64(duration))
		}

		if 200 <= msg.Ctrl.Code && msg.Ctrl.Code < 600 {
			globals.stats.IntStatsInc(fmt.Sprintf("CtrlCodesTotal%dxx", msg.Ctrl.Code/100), 1)
		} else {
			globals.l.Sugar().Warnf("Invalid response code: ", msg.Ctrl.Code)
		}
	}

	select {
	case s.send <- msg:
	default:
		// Never block here since it may also block the topic's run() goroutine.
		globals.l.Sugar().Errorf("s.queueOut: session's send queue full", s.sid)
		return false
	}

	if s.isMultiplex() {
		s.scheduleClusterWriteLoop()
	}

	return true
}
