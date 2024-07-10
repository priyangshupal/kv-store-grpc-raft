package main

import (
	"fmt"
	"time"

	"github.com/priyangshupal/grpc-raft-consensus/store"
)

// sets the configs for a starting raftServer replica
func makeServer(addr string, nodes ...string) *RaftServer {
	raftServerOpts := RaftServerOpts{
		BootstrapNodes: nodes,
		role:           ROLE_FOLLOWER,
		Transport:      &GRPCTransport{ListenAddr: addr},
		logfile:        store.NewLogfile(),
	}
	return NewRaftServer(raftServerOpts)
}
func main() {
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")
	s3 := makeServer(":5000", ":3000", ":4000")
	s4 := makeServer(":6000", ":3000", ":4000", ":5000")
	s5 := makeServer(":7000", ":3000", ":4000", ":5000", ":6000")
	go s1.Start()
	time.Sleep(time.Second * 2)
	go s2.Start()
	time.Sleep(time.Second * 2)
	go s3.Start()
	time.Sleep(time.Second * 2)
	go s4.Start()
	time.Sleep(time.Second * 2)
	go s5.Start()

	time.Sleep(time.Second * 25) // let the system start properly

	// start sending operations from client
	operationNumber := 0
	startTime := time.Now()
	for {
		if time.Since(startTime) > time.Second*1 {
			break
		}
		s1.PerformOperation(fmt.Sprintf("add:%d", operationNumber))
		operationNumber++
	}
	fmt.Printf("time spent: %s, operations: %d\n", time.Since(startTime), operationNumber)
}
