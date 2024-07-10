package store

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-memdb"
	"github.com/priyangshupal/grpc-raft-consensus/schema"
)

// contains methods to interact with a replica's logfile
type Log interface {
	Size() int
	GetFinalIndex() (int, error)
	CommitOperation(int, *Transaction) (int, error)
	ApplyOperation() error
	GetFinalTransaction() (*Transaction, error)
	GetTransactionWithIndex(int) (*Transaction, error)
}

type Transaction struct {
	Index     int
	Operation string
	Term      int
}

type Logfile struct {
	logsCount int
	readyTxn  *Transaction // the last commmitted transaction
	log       *memdb.MemDB
}

func NewLogfile() *Logfile {
	db, err := memdb.NewMemDB(schema.LogfileSchema())
	if err != nil {
		log.Fatalf("error while creating Logfile: %v", err)
	}
	return &Logfile{
		log: db,
	}
}

func (l *Logfile) Size() int { return l.logsCount }

func (l *Logfile) GetFinalIndex() (int, error) {
	txn, err := l.GetFinalTransaction()
	if err != nil {
		return 0, err
	}
	if txn == nil {
		return 0, nil
	}
	return txn.Index, nil
}

// `CommitOperation` is the first step of the two phase commit.
// It is initiated by the `leader` to check whether the requested
// transaction is okay to be committed in the replica
// returns the finalIndex after CommitOperation
func (l *Logfile) CommitOperation(expectedFinalIndex int, txn *Transaction) (int, error) {
	finalIndex, err := l.GetFinalIndex()
	if err != nil {
		return 0, fmt.Errorf("could not fetch final transaction from Logfile: %v", err)
	}
	if finalIndex == expectedFinalIndex {
		// if final index is matching, then add the replica is
		// ready to apply the incoming transaction to the Logfile
		// So, the replica keeps track of this transaction until the
		// second phase of the two phase commit (apply phase)
		l.readyTxn = txn
		return finalIndex, nil
	}
	// otherwise send an error along with `currentFinalIndex`
	return finalIndex, fmt.Errorf("final index (%d) not matching expected final index (%d)", finalIndex, expectedFinalIndex)
}

// `ApplyOperation` is the first step of the two phase commit.
// It is initiated by the `leader` to finally apply the previously
// verified transaction in the `commitOperation` step
func (l *Logfile) ApplyOperation() error {
	// last index will be appended to file
	// prepare the data to commit into logfile
	if l.readyTxn == nil {
		return fmt.Errorf("no transaction ready to apply")
	}
	txn := l.log.Txn(true)
	if err := txn.Insert(schema.TABLE_NAME, l.readyTxn); err != nil {
		return err
	}
	// Apply the transaction to Logfile
	txn.Commit()

	// increase the count of number of rows in Logfile
	// and cache the current final index of the LogFile
	l.logsCount++

	return nil
}
func (l *Logfile) GetFinalTransaction() (*Transaction, error) {
	txn := l.log.Txn(false)
	defer txn.Abort()

	raw, err := txn.Last(schema.TABLE_NAME, "id")
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	return raw.(*Transaction), nil
}
func (l *Logfile) GetTransactionWithIndex(index int) (*Transaction, error) {
	txn := l.log.Txn(false)
	defer txn.Abort()

	raw, err := txn.First(schema.TABLE_NAME, "id", index)
	if err != nil {
		return nil, err
	}
	return raw.(*Transaction), nil
}
func stringifyData(data *Transaction) string {
	return fmt.Sprintf("%d;%s;%d\n", data.Index, data.Operation, data.Term)
}
