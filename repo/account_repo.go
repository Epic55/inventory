package repo

import (
	"context"
	"database/sql"

	"money-transfer/contracts"
	"money-transfer/domain"
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *AccountRepo {
	return &AccountRepo{}
}

func (r *AccountRepo) Retrieve(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM accounts WHERE id=$1", id)
	var (
		bal    int64
		status string
	)
	if err := row.Scan(&id, &bal, &status); err != nil {
		return nil, err
	}
	return domain.NewAccount(id, bal, domain.AccountStatus(status)), nil
}

func (r *AccountRepo) UpdateMut(account *domain.Account) *contracts.Mutation {
	dirty := account.Changes.GetDirty()

	if len(dirty) == 0 {
		return nil
	}

	updates := make(map[string]interface{}, len(dirty))
	for field, value := range dirty {
		updates[field] = value
	}

	return &contracts.Mutation{
		Table:   "accounts",
		ID:      string(account.GetID()),
		Updates: updates,
	}
}
