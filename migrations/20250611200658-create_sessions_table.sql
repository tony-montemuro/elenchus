
-- +migrate Up
CREATE TABLE session (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX session_expiry_idex ON session (expiry);

-- +migrate Down
DROP INDEX session_expiry_idex ON session;
DROP TABLE session;
