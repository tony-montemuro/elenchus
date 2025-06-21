
-- +migrate Up
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idex ON sessions (expiry);

-- +migrate Down
DROP INDEX sessions_expiry_idex ON sessions;
DROP TABLE sessions;
