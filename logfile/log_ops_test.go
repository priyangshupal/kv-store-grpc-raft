package logfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitOperation(t *testing.T) {
	logfile := NewLogfile()

	logfileSize := logfile.Size()
	assert.Equal(t, 0, logfileSize, "Logfile size should be 0 before commit")

	finalIndex, _ := logfile.CommitOperation(0, 0, &Transaction{Index: 1, Operation: "add:3", Term: 1})
	assert.Equal(t, 0, finalIndex, "final index should be 0 for first commit operation")

	logfileSize = logfile.Size()
	assert.Equal(t, 0, logfileSize, "Logfile size should be 0 after commit")
}
func TestApplyOperation(t *testing.T) {
	logfile := NewLogfile()

	finalIndex1, _ := logfile.CommitOperation(0, 0, &Transaction{Index: 1, Operation: "add:3", Term: 1})
	assert.Equal(t, 0, finalIndex1, "final index should be 0 for first commit operation")

	_, err := logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation 1")

	finalIndex2, _ := logfile.CommitOperation(1, 1, &Transaction{Index: 2, Operation: "add:4", Term: 2})
	assert.Equal(t, 1, finalIndex2, "final index should be 1 for second commit operation")

	_, err = logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation 2")

	logfileSize := logfile.Size()
	assert.Equal(t, 2, logfileSize, "Logfile size should be 2")
}
func TestGetTransactionAtIndex(t *testing.T) {
	logfile := NewLogfile()

	finalIndex1, _ := logfile.CommitOperation(0, 0, &Transaction{Index: 1, Operation: "add:3", Term: 1})
	assert.Equal(t, 0, finalIndex1, "final index should be 0 for first commit operation")

	_, err := logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation")

	finalIndex2, _ := logfile.CommitOperation(1, 1, &Transaction{Index: 2, Operation: "add:4", Term: 2})
	assert.Equal(t, 1, finalIndex2, "final index should be 1 for second commit operation")

	_, err = logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation")

	txn, err := logfile.GetTransactionWithIndex(2)
	assert.Equal(t, nil, err, "no error should be present while GetTransactionWithIndex operation")

	// check requested transaction
	assert.Equal(t, &Transaction{Index: 2, Operation: "add:4", Term: 2}, txn, "fetched transaction should match expected")
}
func TestGetLastTransaction(t *testing.T) {
	logfile := NewLogfile()

	finalIndex1, _ := logfile.CommitOperation(0, 0, &Transaction{Index: 1, Operation: "add:3", Term: 1})
	assert.Equal(t, 0, finalIndex1, "final index should be 0 for first commit operation")

	_, err := logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation")

	finalIndex2, _ := logfile.CommitOperation(1, 1, &Transaction{Index: 2, Operation: "add:4", Term: 2})
	assert.Equal(t, 1, finalIndex2, "final index should be 1 for second commit operation")

	_, err = logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation")

	finalIndex3, _ := logfile.CommitOperation(2, 2, &Transaction{Index: 3, Operation: "add:5", Term: 2})
	assert.Equal(t, 2, finalIndex3, "final index should be 2 for third commit operation")

	_, err = logfile.ApplyOperation()
	assert.Equal(t, nil, err, "no error should be present while applying operation")

	txn, err := logfile.GetFinalTransaction()
	assert.Equal(t, nil, err, "no error should be present while GetFinalTransaction operation")

	// check requested transaction
	assert.Equal(t, &Transaction{Index: 3, Operation: "add:5", Term: 2}, txn, "fetched transaction should match expected")
}
func TestStringifyData(t *testing.T) {
	logRow := &Transaction{Index: 1, Operation: "add:3", Term: 2}
	assert.Equal(t, "1;add:3;2\n", stringifyData(logRow))
}
