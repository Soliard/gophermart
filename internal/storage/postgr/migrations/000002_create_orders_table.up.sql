CREATE TABLE orders (
    number VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL,
    accrual DECIMAL(10,2),
    uploaded_at TIMESTAMPTZ NOT NULL
);