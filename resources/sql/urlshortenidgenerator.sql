CREATE TABLE url_id (
    id SERIAL PRIMARY KEY,
    current_id BIGINT NOT NULL,
    range_end BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

