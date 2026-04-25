package transfer_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"money-transfer/contracts"
	"money-transfer/domain"
	"money-transfer/usecases/transfer"
)

type stubAccountRepo struct {
	accounts map[domain.AccountID]*domain.Account
}

func (r *stubAccountRepo) Retrieve(_ context.Context, id domain.AccountID) (*domain.Account, error) {
	acc, ok := r.accounts[id]
	if !ok {
		return nil, fmt.Errorf("account not found: %s", id)
	}
	return acc, nil
}

func (r *stubAccountRepo) UpdateMut(account *domain.Account) *contracts.Mutation {
	dirty := account.Changes.GetDirty()
	if len(dirty) == 0 {
		return nil
	}
	return &contracts.Mutation{
		Table:   "accounts",
		ID:      string(account.GetID()),
		Updates: dirty,
	}
}

func newStub(accounts ...*domain.Account) *stubAccountRepo {
	repo := &stubAccountRepo{accounts: make(map[domain.AccountID]*domain.Account)}
	for _, a := range accounts {
		repo.accounts[a.GetID()] = a
	}
	return repo
}

const (
	idUser1 domain.AccountID = "user1"
	idUser2 domain.AccountID = "user2"
)

func activeAccount(id domain.AccountID, balance int64) *domain.Account {
	return domain.NewAccount(id, balance, domain.AccountStatusActive)
}

func inactiveAccount(id domain.AccountID, balance int64) *domain.Account {
	return domain.NewAccount(id, balance, domain.AccountStatusInactive)
}

func TestExecute_HappyPath(t *testing.T) {
	user1 := activeAccount(idUser1, 100)
	user2 := activeAccount(idUser2, 20)

	uc := transfer.NewInteractor(newStub(user1, user2))
	plan, err := uc.Execute(context.Background(), &transfer.TransferRequest{
		FromAccountID: idUser1,
		ToAccountID:   idUser2,
		Amount:        30,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan == nil {
		t.Fatal("expected a plan, got nil")
	}

	muts := plan.GetMutations()
	if len(muts) != 2 {
		t.Fatalf("expected 2 mutations, got %d", len(muts))
	}

	for _, m := range muts {
		if _, ok := m.Updates["balance"]; !ok {
			t.Errorf("mutation for %s missing balance field", m.ID)
		}
		if len(m.Updates) != 1 {
			t.Errorf("mutation for %s has unexpected extra fields: %v", m.ID, m.Updates)
		}
	}

	if user1.GetBalance() != 70 {
		t.Errorf("user1 balance: want 700, got %d", user1.GetBalance())
	}
	if user2.GetBalance() != 50 {
		t.Errorf("user2 balance: want 500, got %d", user2.GetBalance())
	}
}

func TestExecute_InsufficientFunds(t *testing.T) {
	user1 := activeAccount(idUser1, 5)
	user2 := activeAccount(idUser2, 0)

	uc := transfer.NewInteractor(newStub(user1, user2))
	_, err := uc.Execute(context.Background(), &transfer.TransferRequest{
		FromAccountID: idUser1,
		ToAccountID:   idUser2,
		Amount:        10,
	})

	if !errors.Is(err, domain.ErrInsufficientMoney) {
		t.Fatalf("expected ErrInsufficientFunds, got: %v", err)
	}
}

func TestExecute_InactiveSourceAccount(t *testing.T) {
	user1 := inactiveAccount(idUser1, 100)
	user2 := activeAccount(idUser2, 0)

	uc := transfer.NewInteractor(newStub(user1, user2))
	_, err := uc.Execute(context.Background(), &transfer.TransferRequest{
		FromAccountID: idUser1,
		ToAccountID:   idUser2,
		Amount:        10,
	})

	if !errors.Is(err, domain.ErrInactiveAccount) {
		t.Fatalf("expected ErrInactiveAccount, got: %v", err)
	}
}

func TestExecute_ValidationErrors(t *testing.T) {
	uc := transfer.NewInteractor(newStub())

	cases := []struct {
		name string
		req  *transfer.TransferRequest
	}{
		{"missing from", &transfer.TransferRequest{ToAccountID: idUser2, Amount: 100}},
		{"missing to", &transfer.TransferRequest{FromAccountID: idUser1, Amount: 100}},
		{"same account", &transfer.TransferRequest{FromAccountID: idUser1, ToAccountID: idUser1, Amount: 100}},
		{"0 amount", &transfer.TransferRequest{FromAccountID: idUser1, ToAccountID: idUser2, Amount: 0}},
		{"negative amount", &transfer.TransferRequest{FromAccountID: idUser1, ToAccountID: idUser2, Amount: -50}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tc.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}
