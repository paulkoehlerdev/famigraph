CREATE TABLE IF NOT EXISTS users
(
    handle      TEXT NOT NULL UNIQUE, -- User handle, unique identifier
    credentials TEXT NOT NULL,        -- WebAuthn credentials as a JSON array
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS connections
(
    user1_handle         TEXT NOT NULL,                      -- ID of the first user
    user1_connection_otc TEXT NOT NULL,                      -- Connection OTC
    user2_handle         TEXT NOT NULL,                      -- ID of the second user
    created_at           DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the connection was established
    FOREIGN KEY (user1_handle) REFERENCES users (handle) ON DELETE CASCADE,
    FOREIGN KEY (user2_handle) REFERENCES users (handle) ON DELETE CASCADE,
    UNIQUE (user1_handle, user2_handle),                     -- Ensure unique connections
    UNIQUE (user1_handle, user1_connection_otc)              -- ID should be unique per user to make sure otc are not reused
);

CREATE INDEX IF NOT EXISTS idx_user1_handle ON connections (user1_handle);
CREATE INDEX IF NOT EXISTS idx_user2_handle ON connections (user2_handle);
CREATE INDEX IF NOT EXISTS idx_credentials ON users (handle);
