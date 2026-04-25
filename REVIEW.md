func (uc *Interactor) Execute(ctx context.Context, req *TransferRequest) error {
    source, _ := uc.repo.Retrieve(ctx, req.FromAccountID) ## 1) IT IS BETTER TO DECLARE ERROR AND HANDLE IT CORRECTLY TO PROVIDE CORRECT WORK OF THE METHOD
    dest, _ := uc.repo.Retrieve(ctx, req.ToAccountID) ## 2) IT IS BETTER TO DECLARE ERROR AND HANDLE IT CORRECTLY TO PROVIDE CORRECT WORK OF THE METHOD

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
    if err := uc.db.Apply(mutation2); err != nil { ## 3) NO NEED TO DECLARE err AGAIN, JUST err = ...
        return err
    }

    return nil
}