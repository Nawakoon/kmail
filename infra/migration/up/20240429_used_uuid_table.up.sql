CREATE TABLE IF NOT EXISTS used_uuid (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);