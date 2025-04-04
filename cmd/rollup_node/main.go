package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type RollupNode struct {
	mu       sync.Mutex
	state    map[string]int
	batchLog []string
}

func NewRollupNode() *RollupNode {
	return &RollupNode{
		state:    make(map[string]int),
		batchLog: []string{},
	}
}

func (r *RollupNode) ApplyTransaction(from, to string, amount int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state[from] < amount {
		return fmt.Errorf("insufficient funds")
	}

	r.state[from] -= amount
	r.state[to] += amount
	return nil
}

func (r *RollupNode) GenerateStateRoot() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	data := ""
	for addr, balance := range r.state {
		data += addr + fmt.Sprintf("%d", balance)
	}
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (r *RollupNode) SubmitBatch() {
	stateRoot := r.GenerateStateRoot()
	r.batchLog = append(r.batchLog, stateRoot)
	fmt.Printf("📦 New batch submitted with state root: %s\n", stateRoot)
}

func main() {
	node := NewRollupNode()
	node.state["Alice"] = 100
	node.state["Bob"] = 50

	_ = node.ApplyTransaction("Alice", "Bob", 10)
	_ = node.ApplyTransaction("Bob", "Alice", 5)

	node.SubmitBatch()
}
