CREATE TYPE account_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id),
    name VARCHAR(30) NOT NULL,
    balance INT DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'THB',
    Status account_status DEFAULT 'PENDING'
);