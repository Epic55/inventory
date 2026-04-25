**Q1:** In your implementation, what happens if `source.Withdraw()` succeeds but `dest.Deposit()` fails? Show the exact state of both accounts and the returned plan.

**Q2:** The buggy code applies mutations one at a time. Why is this a problem? Give a specific failure scenario.

**Q3:** Your `UpdateMut` should only include dirty fields. If an account has `balance` changed but `status` unchanged, the mutation should NOT include `status`. Why does this matter for concurrent updates?

**Q4:** Look at this alternative approach:

```go
func (r *AccountRepo) UpdateMut(account *Account) *Mutation {
    return &Mutation{
        Updates: map[string]interface{}{
            "balance": account.Balance(),
            "status":  account.Status(),  // Always include all fields
        },
    }
}
```

What problem does this cause that the dirty-field approach avoids?