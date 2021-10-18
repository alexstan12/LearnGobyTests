package pointersErrors

import (
	"errors"
	"fmt"
)

type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

type Stringer interface {
	String() string
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

func (w *Wallet) Deposit(amount Bitcoin){
	(*w).balance += amount
}

func (w *Wallet) Balance() Bitcoin{
	return w.balance
}

func (w *Wallet) Withdraw(b Bitcoin) error{
	if w.Balance() < b {

		return ErrInsufficientFunds
	}
	w.balance -= b
	return nil
}
