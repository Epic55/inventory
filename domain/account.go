package domain

import "errors"

type AccountID string
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)

var (
	ErrInsufficientMoney = errors.New("Not enough money")
	ErrInactiveAccount   = errors.New("Account is inactive")
	errInvalidAmount     = errors.New("Amount must be greater than 0")
)

type ChangeTracker struct {
	dirty map[string]interface{}
}

type Account struct {
	id      AccountID
	balance int64
	status  AccountStatus
	Changes ChangeTracker
}

func NewAccount(id AccountID, balance int64, status AccountStatus) *Account {
	return &Account{id: id, balance: balance, status: status}
}

func (a *Account) Withdraw(amount int64) error {
	if amount <= 0 {
		return errInvalidAmount
	}
	if a.status != AccountStatusActive {
		return ErrInactiveAccount
	}
	if a.balance < amount {
		return ErrInsufficientMoney
	}
	a.balance -= amount
	a.Changes.Set("balance", a.balance)
	return nil
}

func (a *Account) Deposit(amount int64) error {
	if amount <= 0 {
		return errInvalidAmount
	}
	if a.status != AccountStatusActive {
		return ErrInactiveAccount
	}
	a.balance += amount
	a.Changes.Set("balance", a.balance)
	return nil
}

func (c *ChangeTracker) Set(field string, value interface{}) {
	if c.dirty == nil {
		c.dirty = make(map[string]interface{})
	}
	c.dirty[field] = value
}

func (c *ChangeTracker) GetDirty() map[string]interface{} {
	return c.dirty
}

func (a *Account) GetID() AccountID  { return a.id }
func (a *Account) GetBalance() int64 { return a.balance }
