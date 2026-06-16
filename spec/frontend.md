# Framework

- Nuxt 3 (Vue 3)
- Pinia (State Management)
- Tailwind CSS
- Chart.js or ECharts
- GSAP / Vue Transition for animation effects



# Layout

- Mobile-first responsive design (RWD)
- Vertical stacked layout
- Center-aligned content
- Floating card-based interface
- Generous spacing between sections (16px–24px)
- Desktop version should limit max-width for better visual balance
- Avoid dense grids and enterprise-style layouts


# Others
- 2 languages: Traditional Chinese, English
- Use i18n (tslang)



# Formatting Conventions

## Currency Formatting

- Any monetary value (balance/amount) passed to the currency formatter that is missing, `null`, `undefined`, or non-finite (e.g. `NaN`) must be treated as `0` and displayed as the formatted zero amount (e.g. "$0") — never raw "NaN"/"非數值".


# Components

## pageCard

Used as navigation entry cards.

Features:
- Large rounded card
- Soft hover animation
- Floating feeling
- Icon + title + short description

Used in:
- Home page
- Statistics page



## functionCard

Main functional container.

Features:
- Soft gradient background
- Large spacing
- Rounded corners
- Main interaction area

Used in:
- Booking
- Budget editing
- Setting page



## table

Minimal lightweight table design.

Features:
- No hard borders
- Only subtle separators
- Clean spacing
- Mobile responsive scrolling



# Frontend Pages

## Home (`/`)

### If Logged In

Display 4 navigation cards (`pageCard`):

| Card | Route |
|-|-|
| Set Budget | `/setBudget` |
| Booking | `/booking` |
| Statistics | `/statistics` |
| Setting | `/setting` |

Each card includes:
- Bird-themed icon
- Soft hover animation
- Short description



### If Not Logged In

Display:
- Welcome message
- Small animated bird illustration
- Bird looks confused / thinking
- Floating animated bubbles above bird

Bubble icons:
- Budget ?
- Money ?
- Nest ?
- Egg ?

Below animation:
- Google Login button

Authentication flow:
1. User clicks Google Login
2. Receive Google Token
3. Send token to backend: `POST /user/login`
4. Backend verifies token and returns JWT
5. Store JWT; redirect to home



## Set Budget (`/setBudget`)

### Main View

On mount, call `GET /user/me` to refresh the budgets store with the latest balances/percentages (the in-memory store can be stale after bookings are recorded on other pages or in other tabs).

Display:
- Large pie chart in center showing percentage distribution of all budget groups
- Budget group table (always visible when budgets exist)
- "+ Add Group" button, always visible on the main view



### Group Management

Functions:
- Add group
- Edit group
- Delete group (saving bucket cannot be deleted or edited)

Budget group table columns:
- Icon
- Name (+ optional note)
- Percentage
- Balance (`user_budgets.amount`)
- Actions (edit / delete)

API:
- `POST /setBudget/add`
- `POST /setBudget/change`
- `POST /setBudget/delete`



#### Add Group Form

Top section — new budget fields:
- Group name
- Icon selector (icons already used by the user's other budgets are disabled — data comes from the existing budgets list already loaded on this page)
- Color picker (pastel colors only; colors already used by the user's other budgets are disabled — same source; available colors appear vivid with increased saturation, disabled colors appear faded with low opacity)
- Percentage of income (numeric input, e.g. 30 for 30%)
- Note (optional)

Below form fields — existing budgets list (inline editable):
- All existing budget groups are shown in a list
- Each row displays: icon, name, current percentage (editable input)
- Saving row is shown last with a grayed-out preview of its new percentage:
  - e.g. "Saving: 70% → 50%" (gray text next to the saving row)
  - Updates in real-time as user types the new budget's percentage or edits others
- User can edit any non-saving budget's percentage inline to manually redistribute

Group name validation:
- If the entered name matches an existing budget group's name (case-insensitive), show an inline error ("Name already in use") and disable submit

Submit behavior:
- Submit button is disabled when total ≠ 100%; show a running total indicator (e.g. "Total: 110% — must be 100%")
- Submit button is active only when total == 100% and the name is valid
- On submit: sends new budget fields + all non-saving budgets' updated percentages to `POST /setBudget/add`

#### Edit Group Form

Top section — editable fields for the target budget group (pre-filled with existing values):
- Group name
- Icon selector (icons used by the user's *other* budgets are disabled; the current budget's own icon stays selectable)
- Color picker (pastel colors only; colors used by the user's *other* budgets are disabled; the current budget's own color stays selectable; available colors appear vivid, disabled colors appear faded with low opacity)
- Note (optional)

Below form fields — same inline list of all budgets as Add Group Form (pre-filled with existing percentages), where the user edits percentages directly.

Group name validation:
- If the entered name matches another budget group's name (case-insensitive, excluding this budget's current name), show an inline error ("Name already in use") and disable submit

- Saving row shows grayed-out preview when any percentage is changed
- Submit button disabled until total == 100% and the name is valid
- On submit: sends updated fields (name, icon, color, note) + all non-saving budgets' updated percentages to `POST /setBudget/change`

#### Percentage Validation

- Saving's percentage is never manually editable — it is always computed as 100% - sum(others) and shown as a preview
- Saving cannot go below 0%: if the sum of non-saving budgets would exceed 100%, show an inline error and disable submit
- Backend will also reject if total ≠ 100% (double validation)



## Booking (`/booking`)

### Main Layout

Display one large `functionCard`.

Top section:
- Toggle button: Spending / Income
- Default: Spending selected



### Spending Mode

Fields:
- Budget group dropdown (lists all user_budgets, including saving — users can record spending against their savings)
- Amount
- Note (optional)

Submit: `POST /booking/add`

Feedback:
- Success: Green happy bird animation
- Error: Red warning bird animation



### Income Mode

User chooses distribution method via toggle:

#### Option 1 — Auto Distribution

Automatically distribute income to all budget groups by their percentage.

Fields:
- Amount
- Note (optional)

Submit: `POST /booking/add` (distributed)



#### Option 2 — Specific Budget Group

User selects one target group.

Fields:
- Budget group dropdown (includes all budgets including saving)
- Amount
- Note (optional)

Submit: `POST /booking/add` (not distributed)



### Recent Records Table

Below the form, display the 5 most recent records for the **currently selected tab** (spending or income) for the current month. The table title and content switch when the user toggles between Spending and Income.

- **Spending tab**: shows the 5 most recent spending records; amounts displayed in red with a `-` prefix
- **Income tab**: shows the 5 most recent income records; amounts displayed in green with a `+` prefix

Columns:
- Date
- Type (budget_name; "—" for spending with deleted budget; "Distributed" for auto-distributed income)
- Amount
- Note
- Delete icon

Empty state: show a bird illustration with the appropriate message ("No spending this month yet." / "No income this month yet.")

The table reloads automatically after a successful submit or delete.

API:
- Spending fetch: `GET /report/showSpendingReport?month=&limit=5`
- Income fetch: `GET /report/showIncomeReport?month=&limit=5`
- Delete: `POST /booking/delete`



## Statistics (`/statistics`)

Display 2 navigation cards (`pageCard`):

| Card | Route |
|-|-|
| See Report | `/report` |
| See Each Budget | `/budget` |



## Report (`/report`)

### Filters

Tab 1 — Type:
- Spending (default)
- Income

Dropdown 1 — Scope (only shown in Spending tab):
- All spending (default)
- Specific budget (lists all budget groups)

Dropdown 2 — Period:
- Monthly (default)
- Annual

Dropdown 3 — Time (linked to Dropdown 2):
- When Monthly → year + month picker (e.g. 2026-06)
- When Annual → year picker only (e.g. 2026)

Rule: only show months / years that have existing records.
- API: `GET /report/availablePeriods?type=spending|income` returns the list of months (YYYY-MM) and years that have existing records, used to populate the time picker.
- When Scope = Specific budget, pass `budget_id` to scope the period list to that budget's records.



### Spending Tab

#### Pie Chart
- Scope = All: shows spending distribution across all budget groups
- Scope = Specific budget: hide pie chart; show monthly trend line chart for that budget instead (last 12 months or within the selected year)

#### Summary Table (per budget, Scope = All only)

Columns:
- Budget Group
- Total Spent (for selected period)
- Current Balance (user_budgets.amount, real-time — not historical)

Note: Current Balance always reflects the live balance, not the balance at the time of the selected period.

#### Transaction List (individual records)

Columns:
- Date
- Budget Group (use budget_name if budget was deleted)
- Amount
- Note
- Delete icon → `POST /booking/delete`

Empty state: show bird illustration with "No records found."

API (Spending):
- Monthly all: `GET /report/showSpendingReport?month=`
- Monthly specific budget: `GET /report/showBudgetReport?budget_id=&month=`
- Annual all: `GET /report/showAnnualSpendingReport?year=`
- Annual specific budget: `GET /report/showAnnualBudgetReport?budget_id=&year=`



### Income Tab

#### Transaction List (individual records)

Columns:
- Date
- Type (Distributed / budget_name)
- Amount
- Note
- Delete icon → `POST /booking/delete`

Empty state: show bird illustration with "No records found."

API (Income):
- Monthly: `GET /report/showIncomeReport?month=`
- Annual: `GET /report/showAnnualIncomeReport?year=`



## Budget (`/budget`)

### Budget Selector

Dropdown to choose a specific budget group.



### Display

Show:
- Large title
- Budget icon
- Current balance (user_budgets.amount)

Balance color:
- Positive → Green
- Negative → Red (overdraft)



### Charts

Display monthly spending trend for last 12 months.

Can use:
- Line chart
- Soft mini chart cards

API:
- `GET /report/showAnnualBudgetReport?budget_id=&year=` — call for the current year, and for the previous year too if the last-12-months window spans two calendar years; merge results and slice to the last 12 months



## Setting (`/setting`)

Display one settings card.



### User Information

Read-only:
- Google profile image
- Name
- Email

Cannot be modified.



### Language Setting

Supported languages:
- English
- Traditional Chinese

API: `POST /user/setting`



### Logout

Display a Logout button.

- On click: clear JWT from Pinia store and localStorage, then redirect to `/`
- No backend API call required



# Frontend Global Rules

## Authentication Middleware

- All pages except `/` require a valid JWT
- If not logged in, redirect to `/`
- JWT stored in Pinia store (and persisted to localStorage)
- The budgets list (Pinia `budgets` store) is NOT persisted to localStorage. On app init, if a JWT is restored from localStorage but the budgets store is empty, call `GET /user/me` to rehydrate the user and budgets state before rendering protected pages.

Protected pages:
- `/setBudget`
- `/booking`
- `/statistics`
- `/report`
- `/budget`
- `/setting`

## JWT Handling

- Attach JWT as `Authorization: Bearer <token>` header on every API request
- If backend returns 401, clear JWT and redirect to `/`



## Responsive Design

Must support:
- Mobile phones
- Tablets
- Desktop browsers

Priority: Mobile-first
