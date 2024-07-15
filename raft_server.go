package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/priyangshupal/grpc-raft-consensus/encoding"
	"github.com/priyangshupal/grpc-raft-consensus/fileops"
	"github.com/priyangshupal/grpc-raft-consensus/heartbeat"
	"github.com/priyangshupal/grpc-raft-consensus/logfile"
	"github.com/priyangshupal/grpc-raft-consensus/pb"
	"google.golang.org/grpc"
)

type ROLE int

const (
	ROLE_LEADER    = 1
	ROLE_FOLLOWER  = 2
	ROLE_CANDIDATE = 3
)

const (
	HEARTBEAT_PERIOD  = time.Second * 1
	HEARTBEAT_TIMEOUT = time.Second * 10
)

// type RaftServerOpts struct {
// 	role           ROLE
// 	BootstrapNodes []string
// 	Heartbeat      *heartbeat.Heartbeat
// 	logfile        logs.Log
// 	Transport      Transport
// 	applyCh        chan *logfile.Transaction
// }

type RaftServer struct {
	role           ROLE
	BootstrapNodes []string
	Heartbeat      *heartbeat.Heartbeat
	logfile        logfile.Log
	Transport      Transport
	applyCh        chan *logfile.Transaction

	leaderAddr  string
	currentTerm int
	commitIndex int

	// { address of server, connection client}, we will maintain
	// the connections and reuse them to reduce latency (creating
	// new connections increases latency)
	ReplicaConnMap     map[string]*grpc.ClientConn
	ReplicaConnMapLock sync.RWMutex
}

// Performs the operation requested by the client.
func (s *RaftServer) PerformOperation(operation string) error {
	// only the LEADER is allowed to perform the operation
	// and it then replicates that operation across all the nodes.
	// if the current node is not a LEADER, the operation request
	// will be forwarded to the LEADER, who will then perform the operation
	log.Printf("[%s] received operation (%s)\n", s.Transport.Addr(), operation)
	if s.role == ROLE_LEADER {
		txn, err := s.convertToTransaction(operation)
		if err != nil {
			return fmt.Errorf("[%s] error while converting to transaction", s.Transport.Addr())
		}
		return s.performTwoPhaseCommit(txn)
	}
	log.Printf("[%s] forwarding operation (%s) to leader [%s]\n", s.Transport.Addr(), operation, s.leaderAddr)
	s.ReplicaConnMapLock.RLock()
	defer s.ReplicaConnMapLock.RUnlock()

	// sending operation to the LEADER to perform a TwoPhaseCommit
	return sendOperationToLeader(operation, s.ReplicaConnMap[s.leaderAddr])
}

func (s *RaftServer) convertToTransaction(operation string) (*logfile.Transaction, error) {
	// structure of operation => operationName:value... for eg: "add:5"
	return &logfile.Transaction{Index: s.commitIndex + 1, Operation: operation, Term: s.currentTerm}, nil
}

// `sendOperationToLeader` is called when an operation reaches a FOLLOWER.
// This function forwards the operation to the LEADER
func sendOperationToLeader(operation string, conn *grpc.ClientConn) error {
	replicateOpsClient := pb.NewReplicateOperationServiceClient(conn)
	_, err := replicateOpsClient.ForwardOperation(
		context.Background(),
		&pb.ForwardOperationRequest{Operation: operation},
	)
	if err != nil {
		return err
	}
	return nil
}

// Performs a two phase commit on all the FOLLOWERS
func (s *RaftServer) performTwoPhaseCommit(txn *logfile.Transaction) error {
	s.ReplicaConnMapLock.RLock()
	wg := &sync.WaitGroup{}

	// First phase of the TwoPhaseCommit: Commit operation

	// CommitOperation on self
	if _, err := s.logfile.CommitOperation(s.commitIndex, s.commitIndex, txn); err != nil {
		panic(fmt.Errorf("[%s] %v", s.Transport.Addr(), err))
	}

	log.Printf("[%s] performing commit operation on %d followers\n", s.Transport.Addr(), len(s.ReplicaConnMap))

	for addr, conn := range s.ReplicaConnMap {
		replicateOpsClient := pb.NewReplicateOperationServiceClient(conn)
		log.Printf("[%s] sending (CommitOperation: %s) to [%s]\n", s.Transport.Addr(), txn.Operation, addr)
		response, err := replicateOpsClient.CommitOperation(
			context.Background(),
			&pb.CommitTransaction{
				ExpectedFinalIndex: int64(s.commitIndex),
				Index:              int64(txn.Index),
				Operation:          txn.Operation,
				Term:               int64(txn.Term),
			},
		)
		if err != nil {
			log.Printf("[%s] received error in (CommitOperation) from [%s]: %v", s.Transport.Addr(), addr, err)
			// if there is both, an error and a response, the FOLLOWER is missing
			// some logs. So the LEADER will replicate all the missing logs in the FOLLOWER
			if response != nil {
				wg.Add(1)
				go s.replicateMissingLogs(int(response.LogfileFinalIndex), addr, replicateOpsClient, wg)
			} else {
				return err
			}
		}
	}

	// wait for all FOLLOWERS to be consistent
	wg.Wait()

	log.Printf("[%s] performing (ApplyOperation) on %d followers\n", s.Transport.Addr(), len(s.ReplicaConnMap))

	// Second phase of the TwoPhaseCommit: Apply operation

	// ApplyOperation on self
	if _, err := s.logfile.ApplyOperation(); err != nil {
		panic(err)
	}

	for _, conn := range s.ReplicaConnMap {
		replicateOpsClient := pb.NewReplicateOperationServiceClient(conn)
		_, err := replicateOpsClient.ApplyOperation(
			context.Background(),
			&pb.ApplyOperationRequest{},
		)
		if err != nil {
			return err
		}
	}
	s.ReplicaConnMapLock.RUnlock()

	s.commitIndex++ // increment the final commitIndex after applying changes

	s.applyCh <- txn

	return nil
}

// `replicateMissingLogs` makes a FOLLOWER consistent with the leader. This is
// called when the FOLLOWER is missing some logs and refuses a commit operation
// request from the LEADER
func (s *RaftServer) replicateMissingLogs(startIndex int, addr string, client pb.ReplicateOperationServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		startIndex++
		txn, err := s.logfile.GetTransactionWithIndex(startIndex)
		if err != nil {
			log.Printf("[%s] error fetching (index: %d) from Logfile\n", s.Transport.Addr(), startIndex)
			break
		}
		if txn == nil {
			break
		}
		_, err = client.CommitOperation(
			context.Background(),
			&pb.CommitTransaction{
				ExpectedFinalIndex: int64(startIndex),
				Index:              int64(txn.Index),
				Operation:          txn.Operation,
				Term:               int64(txn.Term),
			},
		)
		if err != nil {
			log.Printf("[%s] error replicating missing log (index: %d) to [%s]\n", s.Transport.Addr(), startIndex, addr)
		}
	}
}

// `requestVotes` is called when the heartbeat has timed out
// and the raft server turns into a candidate.
// It returns the number of votes received along with error (if any)
func (s *RaftServer) requestVotes() int {
	var numVotes int = 0

	s.ReplicaConnMapLock.RLock()
	defer s.ReplicaConnMapLock.RUnlock()

	// iterate over replica addresses and request
	// vote from each replica
	for _, conn := range s.ReplicaConnMap {
		electronServiceClient := pb.NewElectionServiceClient(conn)
		response, err := electronServiceClient.Voting(
			context.Background(),
			&pb.VoteRequest{LogfileIndex: uint64(s.commitIndex)},
		)
		if err != nil {
			log.Printf("[%s] error while requesting vote: %v\n", s.Transport.Addr(), err)
			return 0
		}
		if response.VoteType == pb.VoteResponse_VOTE_GIVEN {
			numVotes += 1
		}
	}
	return numVotes
}

func (s *RaftServer) sendHeartbeat() int {
	aliveCount := 0
	s.ReplicaConnMapLock.RLock()
	for addr, conn := range s.ReplicaConnMap {
		heartbeatClient := pb.NewHeartbeatServiceClient(conn)
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*1))
		defer cancel()
		response, err := heartbeatClient.Heartbeat(
			ctx,
			&pb.HeartbeatRequest{
				IsAlive: true,
				Addr:    s.Transport.Addr(),
			})
		if err != nil {
			log.Printf("[%s] error while sending heartbeat to [%s]: %v\n", s.Transport.Addr(), addr, err)
		}
		if response != nil && response.IsAlive {
			aliveCount++
		}
	}
	s.ReplicaConnMapLock.RUnlock()
	return aliveCount
}

// `sendHeartbeatPeriodically` is called by the leader to
// send a heartbeat to followers every second
func (s *RaftServer) sendHeartbeatPeriodically() {
	// start the process of sending heartbeat for a leader
	for {
		log.Printf("[%s] sending heartbeat to %d followers\n", s.Transport.Addr(), len(s.ReplicaConnMap))
		aliveReplicas := s.sendHeartbeat()
		if aliveReplicas < (len(s.ReplicaConnMap)-1)/2 {
			panic("more than half of the replicas are down")
		}
		time.Sleep(HEARTBEAT_PERIOD)
	}
}

// starts the heartbeat timeout process for a FOLLOWER
func (s *RaftServer) startHeartbeatTimeoutProcess() error {
	log.Printf("[%s] starting heartbeat timeout\n", s.Transport.Addr())
	timeoutFunc := func() {
		log.Printf("[%s] timeout expired\n", s.Transport.Addr())

		s.role = ROLE_CANDIDATE // the replica becomes a CANDIDATE to contest in election
		votesWon := s.requestVotes()
		totalVotes := 1 + votesWon
		totalCandidates := 1 + len(s.ReplicaConnMap)
		// a candidate wins the election and becomes a leader
		// if it receives more than half of the total votes
		if totalVotes >= totalCandidates/2 {
			// if it wins the election, turn it into a LEADER
			// and start sending heartbeat process
			log.Printf("[%s] is the leader", s.Transport.Addr())
			s.role = ROLE_LEADER
			s.currentTerm++
			s.leaderAddr = s.Transport.Addr()

			// if the replica becomes a LEADER, it does not need to listen
			// for heartbeat from other replicas anymore, so stop the
			// heartbeat timeout process
			s.Heartbeat.Stop()

			// the LEADER will send heartbeat to the FOLLOWERS
			go s.sendHeartbeatPeriodically()
		} else {
			// if it loses the election, turn it back in to a follower
			s.role = ROLE_FOLLOWER
		}
	}
	// start/reset heartbeat timeout process for the follower
	// this will trigger the timeoutFunc after a timeout
	if s.role == ROLE_FOLLOWER {
		s.Heartbeat = heartbeat.NewHeartbeat(HEARTBEAT_TIMEOUT, timeoutFunc)
	}
	return nil
}

// sends requests to other replicas so that they can
// add this server to their replicaConnMap
func (s *RaftServer) bootstrapNetwork() {
	wg := &sync.WaitGroup{}
	for _, addr := range s.BootstrapNodes {
		wg.Add(1)
		if len(addr) == 0 {
			continue
		}
		go func(s *RaftServer, addr string, wg *sync.WaitGroup) {
			log.Printf("[%s] attempting to connect with [%s]\n", s.Transport.Addr(), addr)
			if err := s.Transport.Dial(s, addr); err != nil {
				log.Printf("[%s]: dial error while connecting to [%s]: %v\n", s.Transport.Addr(), addr, err)
			}
			wg.Done()
		}(s, addr, wg)
	}
	wg.Wait()
	log.Printf("[%s] bootstrapping completed\n", s.Transport.Addr())
}

func (s *RaftServer) startGrpcServer() error {
	lis, err := net.Listen("tcp", s.Transport.Addr())
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.Transport.Addr(), err)
	}

	// register services with the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterHeartbeatServiceServer(grpcServer, NewHeartbeatServiceServer(s))
	pb.RegisterBootstrapServiceServer(grpcServer, NewBootstrapServiceServer(s))
	pb.RegisterElectionServiceServer(grpcServer, NewElectionServiceServer(s))
	pb.RegisterReplicateOperationServiceServer(grpcServer, NewReplicateOpsServiceServer(s))
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC on port %s: %v", s.Transport.Addr(), err)
	}
	return nil
}

// sets the configs for a starting raftServer replica
func makeRaftServer(addr string, applyCh chan *logfile.Transaction, nodes []string) (*RaftServer, map[string]string) {
	raftServer := &RaftServer{
		BootstrapNodes: nodes,
		role:           ROLE_FOLLOWER,
		Transport:      &GRPCTransport{ListenAddr: addr},
		logfile:        logfile.NewLogfile(),
		applyCh:        applyCh,
		ReplicaConnMap: make(map[string]*grpc.ClientConn),
	}

	// check if a snapshot doesn't exists for the server,
	// then return a new raft server
	filePath := fmt.Sprintf("%s/%s.%s", SNAPSHOTS_DIR, addr, fileops.FILE_EXTENSION)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[%s] snapshot doesn't exist, creating new RaftServer", addr)
		return raftServer, make(map[string]string)
	}
	log.Printf("[%s] restoring RaftServer from snapshot", addr)

	// if a snapshot exists, add the additional configs to server
	// commitIndex, logfilelength, kv map
	snapshotContent, err := fileops.ReadFile(SNAPSHOTS_DIR, addr)
	if err != nil {
		log.Fatalf("error while reading snapshot: %v", err)
	}
	var kvMap map[string]string
	raftServer.commitIndex, kvMap = destructureSnapshot(snapshotContent)
	return raftServer, kvMap
}

func (s *RaftServer) Start() error {
	go s.startGrpcServer()
	time.Sleep(time.Second * 3) // wait for server to start

	log.Printf("[%s] raft server started\n", s.Transport.Addr())

	// send request to bootstrapped servers to
	// add this to replica to their `replicaConnMap`
	s.bootstrapNetwork()

	s.startHeartbeatTimeoutProcess()

	return nil
}

func destructureSnapshot(content []byte) (int, map[string]string) {
	// encodedCommitIndex := bytes.SplitN(content, []byte("\n"), 2)[0]
	encodedKVMap := bytes.SplitN(content, []byte("\n"), 2)[1]
	decodedContent, err := encoding.Decode(content)
	if err != nil {
		log.Fatal("error while decoding commit index", err)
	}
	decodedCommitIndex := strings.Trim(decodedContent, "\n")
	decodedKVMap, err := encoding.DecodeMap(encodedKVMap)
	if err != nil {
		log.Fatal("error while deserializing map:", err)
	}
	commitIndex, _ := strconv.Atoi(decodedCommitIndex)
	return commitIndex, decodedKVMap
}
