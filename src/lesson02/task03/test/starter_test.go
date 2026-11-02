package test

import (
	"testing"

	task "lesson02_task03"
)

// mapStore is a simple in-memory Store.
type mapStore map[string]int

func (m mapStore) Balance(a string) (int, error)    { return m[a], nil }
func (m mapStore) SetBalance(a string, v int) error { m[a] = v; return nil }

// TestStarter is a first check you can run yourself (the green arrow) before you press Check.
// Add your own cases here if you like; Check runs this file together with hidden tests.
func TestStarter(t *testing.T) {
	s := mapStore{"alice": 100, "bob": 0}
	if err := task.Transfer(s, "alice", "bob", 30); err != nil {
		t.Fatalf("Transfer returned %v", err)
	}
	if s["alice"] != 70 || s["bob"] != 30 {
		t.Fatalf("balances after transfer: alice=%d bob=%d, want 70 and 30", s["alice"], s["bob"])
	}
}
