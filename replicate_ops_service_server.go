package main

import (
	"context"
	"log"

	"github.com/priyangshupal/grpc-raft-consensus/logfile"
	"github.com/priyangshupal/grpc-raft-consensus/pb"
)

type ReplicateOpsServiceServer struct {
	pb.UnimplementedReplicateOperationServiceServer
	raftServer *RaftServer
}

func NewReplicateOpsServiceServer(raftServer *RaftServer) *ReplicateOpsServiceServer {
	return &ReplicateOpsServiceServer{raftServer: raftServer}
}

func (s *ReplicateOpsServiceServer) CommitOperation(context context.Context, txn *pb.CommitTransaction) (*pb.CommitOperationResponse, error) {
	log.Printf("[%s] received (CommitOperation: %s) from leader\n", s.raftServer.Transport.Addr(), txn.Operation)
	logfileFinalIndex, err := s.raftServer.logfile.CommitOperation(
		int(txn.ExpectedFinalIndex),
		s.raftServer.commitIndex,
		&logfile.Transaction{Index: int(txn.Index), Operation: txn.Operation, Term: int(txn.Term)},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CommitOperationResponse{LogfileFinalIndex: int64(logfileFinalIndex)}, nil
}

func (s *ReplicateOpsServiceServer) ApplyOperation(context context.Context, txn *pb.ApplyOperationRequest) (*pb.ApplyOperationResponse, error) {
	log.Printf("[%s] received (ApplyOperation) from leader\n", s.raftServer.Transport.Addr())
	appliedTxn, err := s.raftServer.logfile.ApplyOperation()
	if err != nil {
		return nil, err
	}
	s.raftServer.commitIndex++
	s.raftServer.applyCh <- appliedTxn
	return nil, nil
}

func (s *ReplicateOpsServiceServer) ForwardOperation(context context.Context, in *pb.ForwardOperationRequest) (*pb.ForwardOperationResponse, error) {
	txn, err := s.raftServer.convertToTransaction(in.Operation)
	if err != nil {
		return nil, err
	}
	if err = s.raftServer.performTwoPhaseCommit(txn); err != nil {
		return nil, err
	}
	return nil, nil
}
