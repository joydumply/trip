CREATE TABLE IF NOT EXISTS trips (
    id SERIAL PRIMARY KEY,
    leader_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
