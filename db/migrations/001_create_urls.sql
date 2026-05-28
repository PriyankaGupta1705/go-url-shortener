CREATE TABLE IF NOT EXISTS urls (
                                    id          SERIAL PRIMARY KEY,
                                    code        VARCHAR(10) UNIQUE NOT NULL,
    original    TEXT NOT NULL,
    visits      INTEGER DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_urls_code ON urls(code);