# Style: fragment with a starter test

The learner fills in one function. A **visible** starter test in `test/` lets them run a first check before
pressing Check. The **hidden** test checks behaviour instead of output: it passes a recording fake through the
`Store` interface and asserts on what happened (error classes with `errors.Is`, no writes after a failure).

Implement `Transfer`:

1. A non-positive amount returns `ErrInvalidAmount`.
2. Read both balances; return the store's error if an account is missing.
3. If the source cannot cover the amount, return an error that matches `ErrInsufficientFunds`.
4. Only then write both balances. A failed transfer must not change anything.

Use this style for library-style exercises; give the learner a runnable starter test and keep the edge cases hidden.
