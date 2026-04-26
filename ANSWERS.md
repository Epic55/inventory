**Q1:** In your implementation, what happens if `source.Withdraw()` succeeds but `dest.Deposit()` fails? Show the exact state of both accounts and the returned plan.

Money from source account will be withdrawed and will not come to destination account, they dissapear. There will be no returned plan, because 'Execute' method will finish before executing 'UpdateMut' method.

**Q2:** The buggy code applies mutations one at a time. Why is this a problem? Give a specific failure scenario.

If 1st mutation will be applied and 2nd mutation will fail, so money will be withdrawed from 1st account, but not received on 2nd account. It is better to do such operations in 1 transaction to db. If 2nd mutation fails, so 1st mutation will be rolled back.

**Q3:** Your `UpdateMut` should only include dirty fields. If an account has `balance` changed but `status` unchanged, the mutation should NOT include `status`. Why does this matter for concurrent updates?

Because if a request changes balance, we need to be sure that status will not be changed with this request, we ll know that another request changed status without changing balance.

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

This sends every field (including untouched). There is can be problem with concurrent update.