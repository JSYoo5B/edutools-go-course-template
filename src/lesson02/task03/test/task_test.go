package test

import (
	"errors"
	"testing"

	task "lesson02_task03"
)

// recordingStore records writes so the test can check behaviour, not output.
type recordingStore struct {
	balances map[string]int
	writes   int
}

func (r *recordingStore) Balance(a string) (int, error) {
	v, ok := r.balances[a]
	if !ok {
		return 0, errors.New("no such account: " + a)
	}
	return v, nil
}

func (r *recordingStore) SetBalance(a string, v int) error {
	r.writes++
	r.balances[a] = v
	return nil
}

func TestTransferErrorsChangeNothing(t *testing.T) {
	cases := []struct {
		name   string
		amount int
		want   error
	}{
		{"zero amount", 0, task.ErrInvalidAmount},
		{"negative amount", -5, task.ErrInvalidAmount},
		{"more than the balance", 500, task.ErrInsufficientFunds},
	}
	for _, c := range cases {
		s := &recordingStore{balances: map[string]int{"alice": 100, "bob": 0}}
		err := task.Transfer(s, "alice", "bob", c.amount)
		if !errors.Is(err, c.want) {
			t.Errorf("%s: Transfer returned %v, want an error matching %v (errors.Is)", c.name, err, c.want)
		}
		if s.writes != 0 {
			t.Errorf("%s: %d balance writes on a failed transfer, want none", c.name, s.writes)
		}
	}
}

func TestTransferUnknownAccountChangesNothing(t *testing.T) {
	s := &recordingStore{balances: map[string]int{"alice": 100}}
	if err := task.Transfer(s, "alice", "nobody", 10); err == nil {
		t.Fatal("Transfer to an unknown account returned nil, want the store's error")
	}
	if s.writes != 0 || s.balances["alice"] != 100 {
		t.Errorf("a failed transfer changed balances: %v (%d writes)", s.balances, s.writes)
	}
}

func TestTransferWholeBalance(t *testing.T) {
	s := &recordingStore{balances: map[string]int{"alice": 100, "bob": 5}}
	if err := task.Transfer(s, "alice", "bob", 100); err != nil {
		t.Fatalf("moving the whole balance returned %v", err)
	}
	if s.balances["alice"] != 0 || s.balances["bob"] != 105 {
		t.Errorf("balances: %v, want alice=0 bob=105", s.balances)
	}
}
