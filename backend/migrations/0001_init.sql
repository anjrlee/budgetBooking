CREATE TABLE users (
    id         TEXT PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    picture    TEXT,
    language   TEXT NOT NULL DEFAULT 'zh-Hant',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_budgets (
    id         BIGSERIAL PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    note       TEXT,
    icon       TEXT NOT NULL,
    color      TEXT NOT NULL,
    percentage NUMERIC(5, 2) NOT NULL DEFAULT 0,
    amount     NUMERIC(14, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, name)
);

CREATE TABLE budget_bookings (
    id          BIGSERIAL PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    budget_id   BIGINT REFERENCES user_budgets(id) ON DELETE SET NULL,
    budget_name TEXT,
    is_income   BOOLEAN NOT NULL,
    amount      NUMERIC(14, 2) NOT NULL,
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_budgets_user_id ON user_budgets(user_id);
CREATE INDEX idx_budget_bookings_user_id ON budget_bookings(user_id);
CREATE INDEX idx_budget_bookings_user_id_created_at ON budget_bookings(user_id, created_at);
