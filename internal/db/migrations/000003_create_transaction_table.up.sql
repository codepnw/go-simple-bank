CREATE TYPE transaction_type AS ENUM ('DEPOSIT', 'TRANSFER', 'WITHDRAW');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_account INT REFERENCES accounts(id),
    to_account INT REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    type transaction_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);