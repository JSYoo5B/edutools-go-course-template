// Package transfer moves money between accounts kept in a Store.
package transfer

import (
	"errors"
	"fmt"
)

// Store keeps account balances.
type Store interface {
	Balance(account string) (int, error)
	SetBalance(account string, amount int) error
}

var (
	// ErrInvalidAmount is returned for an amount that is not positive.
	ErrInvalidAmount = errors.New("amount must be positive")
	// ErrInsufficientFunds is returned when the source account cannot cover the amount.
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Transfer moves amount from one account to another. On any error no balance is changed.
func Transfer(s Store, from, to string, amount int) error {
	//@@begin
	//@@ph // TODO: validate amount, read both balances, check funds, then write both balances.
	//@@ph // Wrap or return ErrInvalidAmount / ErrInsufficientFunds so callers can use errors.Is.
	//@@ph _ = fmt.Errorf
	//@@ph return nil
	if amount <= 0 {
		return ErrInvalidAmount
	}
	fromBal, err := s.Balance(from)
	if err != nil {
		return err
	}
	toBal, err := s.Balance(to)
	if err != nil {
		return err
	}
	if fromBal < amount {
		return fmt.Errorf("transfer %d from %s (balance %d): %w", amount, from, fromBal, ErrInsufficientFunds)
	}
	if err := s.SetBalance(from, fromBal-amount); err != nil {
		return err
	}
	return s.SetBalance(to, toBal+amount)
	//@@end
}
