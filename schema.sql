CREATE TABLE urls (
                      id SERIAL PRIMARY KEY,
                      original_url TEXT NOT NULL UNIQUE,
                      alias VARCHAR(6) NOT NULL UNIQUE,
                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                      click_count INTEGER DEFAULT 0
);