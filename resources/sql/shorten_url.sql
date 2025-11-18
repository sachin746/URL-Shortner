CREATE TABLE shortened_url (
    id PRIMARY KEY NOT NULL, 
    original_url TEXT NOT NULL,
    short_code VARCHAR(12) UNIQUE NOT NULL,
    user_id VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    valid_till TIMESTAMP WITH TIME ZONE
);
