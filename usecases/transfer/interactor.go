package transfer

import (
	"context"
	"errors"
	"fmt"

	"money-transfer/contracts"
	"money-transfer/domain"
)

var (
	errMissingFromAccount = errors.New("FromAccountID is required")
	errMissingToAccount   = errors.New("ToAccountID is required")
	errSameAccount        = errors.New("Source and destination accounts must differ")
	errInvalidAmount      = errors.New("Amount must be greater than 0")
)

type TransferRequest struct {
	FromAccountID domain.AccountID
	ToAccountID   domain.AccountID
	Amount        int64
}

type Interactor struct {
	accounts contracts.AccountRepository
}

func NewInteractor(accounts contracts.AccountRepository) *Interactor {
	return &Interactor{accounts: accounts}
}

func (uc *Interactor) Execute(ctx context.Context, req *TransferRequest) (*contracts.Plan, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	from, err := uc.accounts.Retrieve(ctx, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("Retrieve source account: %w", err)
	}

	to, err := uc.accounts.Retrieve(ctx, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("Retrieve destination account: %w", err)
	}

	if err = from.Withdraw(req.Amount); err != nil {
		return nil, err
	}

	if err = to.Deposit(req.Amount); err != nil {
		return nil, err
	}

	plan := contracts.NewPlan()
	plan.Add(uc.accounts.UpdateMut(from))
	plan.Add(uc.accounts.UpdateMut(to))

	return plan, nil
}

func validateRequest(req *TransferRequest) error {
	if req.FromAccountID == "" {
		return errMissingFromAccount
	}
	if req.ToAccountID == "" {
		return errMissingToAccount
	}
	if req.FromAccountID == req.ToAccountID {
		return errSameAccount
	}
	if req.Amount <= 0 {
		return errInvalidAmount
	}
	return nil
}
