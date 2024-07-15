package main

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/priyangshupal/grpc-raft-consensus/encoding"
	"github.com/priyangshupal/grpc-raft-consensus/fileops"
	"github.com/priyangshupal/grpc-raft-consensus/logfile"
)

const (
	OPERATION_GET    = "GET"
	OPERATION_PUT    = "PUT"
	OPERATION_APPEND = "APPEND"
)

const SNAPSHOTS_DIR = "snapshots/"
const MAX_RAFT_STATE = 100

type KVServer struct {
	rf      *RaftServer
	applyCh chan *logfile.Transaction

	maxraftstate int // snapshot if log grows this long

	kv     map[string]string // { key, Transaction }, operation => opn:3,2
	kvLock sync.RWMutex
}

func NewKVServer() *KVServer {
	return &KVServer{
		kv: make(map[string]string),
	}
}

// `ApplyOperation` function will be called by the client, to make
// any changes to the raft cluster. This will
func (s *KVServer) ApplyOperation(operation string) (string, error) {
	op_type := strings.Split(operation, ":")[0]
	if strings.Compare(op_type, OPERATION_GET) == 0 {
		key := strings.Split(operation, ":")[1]
		return s.get(key), nil
	}
	return "", s.rf.PerformOperation(operation)
}

// background service to listen for applyOperation from Raft server
func (s *KVServer) applyTransactionLoop() {
	for {
		txn := <-s.applyCh
		if txn == nil {
			continue
		}
		// received transaction from Raft server, now replicate
		// the operation in KV store
		s.applyOperationToKVStore(txn)
	}
}

func (s *KVServer) applyOperationToKVStore(txn *logfile.Transaction) {
	log.Printf("[%s] applying operation (%s) to KV store\n", s.rf.Transport.Addr(), txn.Operation)

	operation := txn.Operation
	op_type := strings.Split(operation, ":")[0]
	kv_pair := strings.Split(operation, ":")[1]
	key := strings.Split(kv_pair, ",")[0]
	value := strings.Split(kv_pair, ",")[1]

	switch op_type {
	case OPERATION_PUT:
		s.put(key, value)
	case OPERATION_APPEND:
		s.append(key, value)
	}
	// create a snapshot if the logfile length reaches threshold
	if s.rf.logfile.Size() == s.maxraftstate {
		go s.saveSnapshot()
	}
}

func (s *KVServer) saveSnapshot() {
	// first save the key value store on disk
	log.Printf("[%s] saving snapshot\n", s.rf.Transport.Addr())

	fileName := s.rf.Transport.Addr()

	// write the commit index to the file
	encodedCommitIndex, _ := encoding.Encode(strconv.Itoa(s.rf.commitIndex) + "\n")
	if err := fileops.WriteToFile(SNAPSHOTS_DIR, fileName, encodedCommitIndex); err != nil {
		log.Println("error while writing commit index to file")
		return
	}

	// append the encoded KV map to the same file
	encodedKVMap, err := encoding.Encode(s.kv)
	if err != nil {
		log.Println("error while encoding key/value store")
		return
	}
	if err = fileops.AppendToFile(SNAPSHOTS_DIR, fileName, encodedKVMap); err != nil {
		log.Println("error while writing key/value store to file")
		return
	}

	// then ask the raft server to delete the snapshotted
	// entries from logfile
	_, err = s.rf.logfile.RemoveEntries(s.rf.commitIndex)
	if err != nil {
		log.Printf("[%s] error while deleting entries from logfile: %v", s.rf.leaderAddr, err)
		return
	}
}

func (s *KVServer) put(key string, value string) {
	s.kvLock.Lock()
	defer s.kvLock.Unlock()
	s.kv[key] = value
}

func (s *KVServer) get(key string) string {
	s.kvLock.RLock()
	defer s.kvLock.RUnlock()
	if val, keyExists := s.kv[key]; keyExists {
		return val
	}
	return ""
}

func (s *KVServer) append(key string, value string) {
	s.kvLock.Lock()
	defer s.kvLock.Unlock()
	s.kv[key] += value
}

func StartKVServer(addr string, peers []string) *KVServer {
	kvServer := NewKVServer()
	kvServer.applyCh = make(chan *logfile.Transaction)

	kvServer.maxraftstate = MAX_RAFT_STATE

	kvServer.rf, kvServer.kv = makeRaftServer(addr, kvServer.applyCh, peers)
	kvServer.rf.Start()

	go kvServer.applyTransactionLoop()

	return kvServer
}
