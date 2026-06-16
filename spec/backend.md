# Framework

Language:
- Go (Golang)

Recommended framework:
- Gin
or
- Echo



# Authentication

Use:
- Google OAuth2
- JWT Session

Flow:
1. Frontend receives Google token
2. Backend validates token
3. Create or retrieve user
4. Return JWT token


# Service
## /budget
- /addBudget(userid, name, note, icon, color, percentage)
    - write into db: user_budgets
    - amount default=0

- /syncPercentages(userId, updatedPercentages map[budgetId → percentage])
    - receives a map of all non-saving budgets' new percentages
    - computes saving's new percentage = 100% - sum(updatedPercentages)
    - updates all affected user_budgets.percentage in db
    - validates total == 100%; returns error if not
    - must be called inside the caller's DB transaction (not a standalone transaction)



# API Design

## Auth
- All endpoints except `/user/login` require a valid JWT in the `Authorization: Bearer <token>` header.
- All queries and mutations are automatically scoped to the authenticated user (`user_id` extracted from JWT). The backend must never return or modify another user's data.

## /user

- /login(POST)
    - Google login
    - if user already exists (in db table: users), return information from db
    - if not in db, execute in a single DB transaction:
        - create user record (language default: zh)
        - create saving bucket via addBudget(in service):
            - name='saving', percentage=100%, icon='savings', color='#B5EAD7'
            - this is always the first budget created for the user
        - return user info (do not re-query, use the records just created)

- /me(GET)
    - Requires JWT
    - Returns `{ user, budgets }` for the authenticated user (same shape as the `/user/login` response, minus `token`)
    - Used by the frontend to rehydrate the budgets store after a page refresh / session restore

- /setting(POST)
    - Update `language` field in db table users (only field modifiable via this endpoint)



## /setBudget

- /add(POST)
    - frontend sends new budget fields + all non-saving budgets' new percentages (map of budgetId → percentage)
    - execute in a single DB transaction:
        - insert new budget via addBudget(in service)
        - call syncPercentages(userId, updatedPercentages) to update all percentages including saving
- /change(POST)
    - if user_budget.name == 'saving', reject (saving cannot be directly modified by user)
    - modifiable fields: name, note, icon, color, percentage
    - frontend sends updated fields + all non-saving budgets' new percentages (only required when percentage is being changed)
    - execute in a single DB transaction:
        - update target budget's non-percentage fields in user_budgets
        - if percentage is being changed: call syncPercentages(userId, updatedPercentages)
- /delete(POST)
    - if user_budget.name == 'saving', reject (cannot be deleted)
    - execute in a single DB transaction:
        - add user_budgets.amount (current balance) into the saving bucket's amount
        - add user_budgets.percentage into the saving bucket's percentage
        - delete the budget from user_budgets (triggers ON DELETE SET NULL on budget_bookings)

## /report

Monthly endpoints:
- /showSpendingReport(GET)
    - query params: month (e.g. 2026-06), limit (optional, e.g. 5)
    - query budget_bookings for that month where is_income==false, ordered by created_at desc
    - if limit is provided, return at most that many records
- /showBudgetReport(GET)
    - query params: budget_id, month (e.g. 2026-06)
    - query budget_bookings for that month where is_income==false and budget_id=given id
- /showIncomeReport(GET)
    - query params: month (e.g. 2026-06), limit (optional, e.g. 10)
    - query budget_bookings for that month where is_income==true, ordered by created_at desc
    - if limit is provided, return at most that many records

Period lookup:
- /availablePeriods(GET)
    - query params: type (spending | income), budget_id (optional)
    - returns the distinct months (YYYY-MM) and years that have existing budget_bookings records matching is_income for that type
    - if budget_id is provided, restrict to records for that budget only (used by Spending tab's "Specific budget" scope)
    - used by the frontend to populate the month/year picker (only show periods with existing records)

Annual endpoints:
- /showAnnualSpendingReport(GET)
    - query params: year (e.g. 2026)
    - query budget_bookings for that year where is_income==false
- /showAnnualBudgetReport(GET)
    - query params: budget_id, year (e.g. 2026)
    - query budget_bookings for that year where is_income==false and budget_id=given id
- /showAnnualIncomeReport(GET)
    - query params: year (e.g. 2026)
    - query budget_bookings for that year where is_income==true

## /booking
- /add(POST)
    - execute in a single DB transaction:
        - if spending: insert budget_bookings (budget_id, budget_name snapshot, is_income=false, amount); deduct amount from that budget's user_budgets.amount (overdraft allowed)
        - if income and not distributed: insert budget_bookings (budget_id, budget_name snapshot, is_income=true, amount); add amount to that budget's user_budgets.amount
        - if income and distributed: insert budget_bookings (budget_id=NULL, budget_name=NULL, is_income=true, amount); distribute amount to each user_budgets.amount proportionally by user_budgets.percentage; any rounding remainder goes to saving
- /delete(POST)
    - execute in a single DB transaction:
        - delete the record from budget_bookings
        - reverse the effect on user_budgets.amount:
            - if spending:
                - if budget still exists (budget_id IS NOT NULL): add amount back to that budget's user_budgets.amount
                - if budget was deleted (budget_id IS NULL, budget_name IS NOT NULL): add amount back to saving bucket's amount
            - if income and not distributed:
                - if budget still exists (budget_id IS NOT NULL): deduct amount from that budget's user_budgets.amount
                - if budget was deleted (budget_id IS NULL, budget_name IS NOT NULL): deduct amount from saving bucket's amount
            - if income and distributed (budget_id IS NULL, budget_name IS NULL): reverse-distribute using **current** user_budgets.percentage values
                - note: if percentages have changed since the original booking, the reversal may be slightly inaccurate; the user should be informed of this limitation
