func (uc *Interactor) Execute(ctx context.Context, req *TransferRequest) error {
## 1) need to validate 'req' before doing operations below on it
## 2) it is better to declare error and handle it correctly to provide error handling in lines 5 and 6
## 3) 'uc' variable has 'accounts' field, not 'repo', correct use - uc.accounts.Retrieve in lines 5 and 6
    source, _ := uc.repo.Retrieve(ctx, req.FromAccountID) 
    dest, _ := uc.repo.Retrieve(ctx, req.ToAccountID) 

    source.balance -= req.Amount
    dest.balance += req.Amount

    mutation1 := &Mutation{
        Table:   "accounts",
        ID:      string(source.id),
        Updates: map[string]interface{}{"balance": source.balance},
    }

    mutation2 := &Mutation{
        Table:   "accounts",
        ID:      string(dest.id),
        Updates: map[string]interface{}{"balance": dest.balance},
    }

    if err := uc.db.Apply(mutation1); err != nil {
        return err
    }
## 4) NO NEED TO DECLARE err AGAIN, JUST err = ...
    if err := uc.db.Apply(mutation2); err != nil { 
        return err
    }

    return nil
}