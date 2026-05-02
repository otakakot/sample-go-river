CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
