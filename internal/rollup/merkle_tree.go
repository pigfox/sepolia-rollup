package rollup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

type MerkleTree struct {
	Leaves []string
	Root   string
}

func NewMerkleTree(leaves []string) *MerkleTree {
	mt := &MerkleTree{Leaves: leaves}
	mt.Root = mt.computeRoot()
	return mt
}

func (m *MerkleTree) computeRoot() string {
	if len(m.Leaves) == 0 {
		return ""
	}

	current := m.Leaves
	for len(current) > 1 {
		var next []string
		for i := 0; i < len(current); i += 2 {
			if i+1 == len(current) {
				next = append(next, hashPair(current[i], current[i]))
			} else {
				next = append(next, hashPair(current[i], current[i+1]))
			}
		}
		current = next
	}
	return current[0]
}

func hashPair(a, b string) string {
	sorted := []string{a, b}
	sort.Strings(sorted)
	h := sha256.New()
	h.Write([]byte(sorted[0] + sorted[1]))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MerkleTree) GetProof(index int) []string {
	var proof []string
	current := m.Leaves
	i := index
	for len(current) > 1 {
		var next []string
		for j := 0; j < len(current); j += 2 {
			if j+1 == len(current) {
				next = append(next, hashPair(current[j], current[j]))
			} else {
				next = append(next, hashPair(current[j], current[j+1]))
			}
			if j == i || j+1 == i {
				sibling := j
				if j == i {
					sibling = j + 1
				} else {
					sibling = j
				}
				if sibling < len(current) {
					proof = append(proof, current[sibling])
				}
				i = j / 2
			}
		}
		current = next
	}
	return proof
}

func VerifyProof(leaf string, proof []string, root string) bool {
	h := leaf
	for _, sibling := range proof {
		h = hashPair(h, sibling)
	}
	return h == root
}

// === Example usage ===
func GenerateStateRootFromMap(state map[string]int) (string, []string) {
	var leaves []string
	for k, v := range state {
		leaves = append(leaves, fmt.Sprintf("%s:%d", k, v))
	}
	tree := NewMerkleTree(leaves)
	return tree.Root, leaves
}
