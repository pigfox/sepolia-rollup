package rollup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type RollupNode struct {
	mu         sync.Mutex
	state      map[string]int
	stateRoots []string
	batches    [][]Transaction
}

type Transaction struct {
	From  string
	To    string
	Value int
}

func NewRollupNode() *RollupNode {
	return &RollupNode{
		state: make(map[string]int),
	}
}

func (r *RollupNode) ApplyTransaction(from, to string, value int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state[from] < value {
		return fmt.Errorf("insufficient funds: %s has %d", from, r.state[from])
	}

	r.state[from] -= value
	r.state[to] += value

	tx := Transaction{From: from, To: to, Value: value}
	if len(r.batches) == 0 {
		r.batches = append(r.batches, []Transaction{})
	}
	r.batches[len(r.batches)-1] = append(r.batches[len(r.batches)-1], tx)

	return nil
}

func (r *RollupNode) SubmitBatch() {
	r.mu.Lock()
	defer r.mu.Unlock()

	stateRoot := r.computeStateRoot()
	r.stateRoots = append(r.stateRoots, stateRoot)

	fmt.Printf("🧾 Submitted Batch | New State Root: %s\n", stateRoot)

	// Start a new batch
	r.batches = append(r.batches, []Transaction{})
}

func (r *RollupNode) computeStateRoot() string {
	// Create a deterministic state string
	stateString := ""
	for k, v := range r.state {
		stateString += fmt.Sprintf("%s:%d|", k, v)
	}
	hash := sha256.Sum256([]byte(stateString))
	return hex.EncodeToString(hash[:])
}
