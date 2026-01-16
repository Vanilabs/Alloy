-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN employee_number varchar(50) UNIQUE;
CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS employee_number;
DROP INDEX IF EXISTS users_email_idx;
-- +goose StatementEnd
