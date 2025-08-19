CREATE TABLE withdrawals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    sum DECIMAL(10,2) NOT NULL CHECK (sum >= 0),
    processed_at TIMESTAMPTZ NOT NULL
)