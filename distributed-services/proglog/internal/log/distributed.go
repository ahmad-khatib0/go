package log

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	api "github.com/Ahmadkhatib0/go/distributed-services/proglog/api/v1"
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

///////////////////////////////// LOG api ////////////////////////////////////////
// The DistributedLog will have the same API as the Log type to make them interchangeable

// Append(record *api.Record) appends the record to the log
// we tell Raft to apply a command (we’ve reused for the ProduceRequest for the command)
// that tells the FSM to append the record to the log. Raft runs the process to replicate the command to a
// majority of the Raft servers and ultimately append the record to a majority of Raft servers
func (l *DistributedLog) Append(record *api.Record) (uint64, error) {
	res, err := l.apply(AppendRequestType, &api.ProduceRequest{Record: record})
	if err != nil {
		return 0, err
	}

	return res.(*api.ProduceResponse).Offset, nil
}

// apply(reqType RequestType, req proto.Marshaler) wraps Raft’s API to apply requests and return their responses.
func (l *DistributedLog) apply(reqType RequestType, req proto.Message) (interface{}, error) {
	var buf bytes.Buffer

	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}

	timeout := 10 * time.Second
	future := l.raft.Apply(buf.Bytes(), timeout)
	//  future.Error() API returns an error when something went wrong with Raft’s replication
	if future.Error() != nil {
		return nil, future.Error()
	}

	res := future.Response()
	// opposed to Go’s convention of using Go’s multiple return values to separate errors,
	// you must return a single value for Raft
	if err, ok := res.(error); ok { // assert if its an error
		return nil, err
	}

	return res, nil
}

// Read(offset uint64) reads the record for the offset from the server’s log
func (l *DistributedLog) Read(offset uint64) (*api.Record, error) {
	return l.log.Read(offset)

}

var _ raft.FSM = (*fsm)(nil)

// fsm finite-state-machine
type fsm struct {
	log *Log
}

var _ raft.LogStore = (*logStore)(nil)

// We’re using our own log as Raft’s log store
type logStore struct {
	*Log
}

type RequestType uint8

const (
	AppendRequestType RequestType = 0
)

func newLogStore(dir string, c Config) (*logStore, error) {
	log, err := NewLog(dir, c)
	if err != nil {
		return nil, err
	}
	return &logStore{log}, nil
}

func (l *logStore) FirstIndex() (uint64, error) {
	return l.LowestOffset()
}

func (l *logStore) LastIndex() (uint64, error) {
	off, err := l.HighestOffset()
	return off, err
}

func (l *logStore) GetLog(index uint64, out *raft.Log) error {
	in, err := l.Read(index)
	if err != nil {
		return err
	}

	out.Data = in.Value              // rpc value
	out.Index = in.Offset            // value offset  ( NOTE: WHAT WE CALL OFFSETS, RAFT CALLS INDEXES.)
	out.Type = raft.LogType(in.Type) // value type
	out.Term = in.Term               //
	return nil
}

func (l *logStore) StoreLog(record *raft.Log) error {
	return l.StoreLogs([]*raft.Log{record})
}

func (l *logStore) StoreLogs(records []*raft.Log) error {
	for _, record := range records {
		if _, err := l.Append(&api.Record{
			Value: record.Data,
			Term:  record.Term,
			Type:  uint32(record.Type),
		}); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRange(min, max uint64) a method to delete old records
func (l *logStore) DeleteRange(min, max uint64) error {
	return l.Truncate(max) // removes all segments whose highest offset is lower than lowest.
}

func (l *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])

	switch reqType {
	case AppendRequestType:
		return l.applyAppend(buf[1:])
	}

	return nil
}

func (l *fsm) applyAppend(b []byte) interface{} {
	var req api.ProduceRequest

	err := proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	offset, err := l.log.Append(req.Record)
	if err != nil {
		return err
	}

	return &api.ProduceResponse{Offset: offset}
}

// Snapshot() returns an FSMSnapshot that represents a point-in-time snapshot of the FSM’s state
// These snapshots serve two purposes: they allow Raft to compact its log so it
// doesn’t store logs whose commands Raft has applied already. And they allow Raft to bootstrap new
// servers more efficiently than if the leader had to replicate its entire log again and again.
// /
// Raft calls Snapshot() according to your configured
// SnapshotInterval (how often Raft checks if it should snapshot—default is two minutes) and
// SnapshotThreshold  (how many logs since the last snapshot before making a new snapshot—default is 8192).
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	r := f.log.Reader()

	// call Reader() to return an io.Reader that will read all the log’s data
	return &snapshot{reader: r}, nil
}

var _ raft.FSMSnapshot = (*snapshot)(nil)

type snapshot struct {
	reader io.Reader
}

// Raft calls Persist() on the FSMSnapshot we created to write its state to some sink that, depending on
// the snapshot store we configured Raft with, could be in-memory, a file, an S3 bucket—something
// to store the bytes in.
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	if _, err := io.Copy(sink, s.reader); err != nil {
		_ = sink.Cancel()
		return err
	}
	return sink.Close()
}

// Raft calls Release() when it’s finished with the snapshot.
func (s *snapshot) Release() {}

// Raft calls Restore() to restore an FSM from a snapshot
// In our Restore() implementation, we reset the log and configure its initial offset
// to the first record’s offset we read from the snapshot so the log’s offsets match.
// Then we read the records in the snapshot and append them to our new log.
func (f *fsm) Restore(r io.ReadCloser) error {
	b := make([]byte, lenWidth)
	var buf bytes.Buffer

	for i := 0; ; i++ {
		_, err := io.ReadFull(r, b)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		size := int64(enc.Uint64(b))
		if _, err = io.CopyN(&buf, r, size); err != nil {
			return err
		}

		record := &api.Record{}
		if err = proto.Unmarshal(buf.Bytes(), record); err != nil {
			return err
		}

		if i == 0 {
			f.log.Config.Segment.InitialOffset = record.Offset
			if err := f.log.Reset(); err != nil {
				return err
			}
		}

		if _, err = f.log.Append(record); err != nil {
			return err
		}
		buf.Reset()
	}
	return nil
}