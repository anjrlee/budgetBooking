# Database Engine

- PostgreSQL



# Database Tables

## `users`

| Column | Type | Description |
|-|-|-|
| id | UUID | Primary Key |
| email | Text | Google Email, UNIQUE |
| name | Text | Display Name, NOT NULL |
| language | Text | en / zh |
| created_at | Timestamp | Join Date |

## `user_budgets`

| Column | Type | Description |
|-|-|-|
| id | Int | Primary Key (auto-increment) |
| user_id | UUID | Foreign Key → users.id |
| name | Text | Budget Name, NOT NULL |
| note | Text | Budget Note, nullable |
| icon | Text | Icon ID, NOT NULL |
| color | Text | Hex Color, NOT NULL |
| percentage | NUMERIC(5,2) | Income distribution percentage (0.00–100.00), NOT NULL |
| amount | NUMERIC(15,2) | Current balance, default=0. Can be negative (overdraft allowed) |

Notes:
- All budgets for a user must have percentages that sum to 100%.
- The saving bucket is identified by `name = 'saving'`. It cannot be deleted or modified.
- The saving bucket is always the first budget created per user (at registration).
- `UNIQUE(user_id, name)`: a user cannot have two budgets with the same name.
- `UNIQUE(user_id, icon)`: a user cannot have two budgets with the same icon.
- `UNIQUE(user_id, color)`: a user cannot have two budgets with the same color.



## `budget_bookings`

| Column | Type | Description |
|-|-|-|
| id | BigInt | Primary Key (auto-increment) |
| budget_id | Int | Nullable FK → user_budgets.id ON DELETE SET NULL |
| budget_name | Text | Snapshot of user_budgets.name at booking time, nullable |
| user_id | UUID | Foreign Key → users.id |
| is_income | Boolean | True = Income |
| amount | NUMERIC(15,2) | Transaction Amount |
| note | Text | User Note, nullable |
| created_at | Timestamp | Transaction Time |

Notes:
- `is_income = false` means spending record.
- `amount` is always positive. Sign is determined by `is_income`.
- Use `budget_id` + `budget_name` together to determine the state of a record:

  | budget_id | budget_name | Meaning |
  |-----------|-------------|---------|
  | has value | has value   | normal record, budget still exists |
  | NULL      | has value   | budget was deleted; name is preserved for display |
  | NULL      | NULL        | distributed income record |

- When a `user_budgets` row is deleted, all related `budget_bookings.budget_id` are set to NULL (ON DELETE SET NULL); `budget_name` is retained.
- When deleting a distributed income record, the reversal uses the **current** `user_budgets.percentage` values. This is an accepted limitation — the user should be notified if percentages have changed since the original booking.
- Rounding: when distributing income proportionally, any leftover cent (due to rounding) is assigned to the saving bucket.
