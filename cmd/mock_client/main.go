package main

import (
	"fmt"
	"sepolia-rollup/internal/rollup"
	"time"
)

func main() {
	fmt.Println("🚀 Mock Rollup Client Started")

	node := rollup.NewRollupNode()

	go func() {
		for i := 0; i < 5; i++ {
			err := node.ApplyTransaction("Alice", "Bob", 5)
			if err == nil {
				fmt.Printf("🔄 Alice → Bob: 5 tokens (Txn %d)\n", i+1)
			} else {
				fmt.Println("❌ Transaction failed:", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			node.SubmitBatch()
		}
	}()

	select {} // Keep running
}
