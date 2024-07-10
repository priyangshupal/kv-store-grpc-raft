package main

import (
	"context"
	"log"

	"github.com/priyangshupal/grpc-raft-consensus/pb"
	"github.com/priyangshupal/grpc-raft-consensus/store"
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
		int(txn.ExpectedPreviousIndex),
		&store.Transaction{Index: int(txn.Index), Operation: txn.Operation, Term: int(txn.Term)},
	)
	if err != nil {
		return nil, err
	}
	return &pb.CommitOperationResponse{LogfileFinalIndex: int64(logfileFinalIndex)}, nil
}

func (s *ReplicateOpsServiceServer) ApplyOperation(context context.Context, txn *pb.ApplyOperationRequest) (*pb.ApplyOperationResponse, error) {
	log.Printf("[%s] received (ApplyOperation) from leader\n", s.raftServer.Transport.Addr())
	if err := s.raftServer.logfile.ApplyOperation(); err != nil {
		return nil, err
	}
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
