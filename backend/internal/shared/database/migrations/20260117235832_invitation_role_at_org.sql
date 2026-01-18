-- +goose Up
-- +goose StatementBegin
ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS role_at_org VARCHAR(255) NOT NULL;
ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS department VARCHAR(255) NOT NULL;
ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS phone VARCHAR(255) NOT NULL;
ALTER TABLE invitations
    ADD CONSTRAINT invitations_status_check CHECK (status IN ('accepted', 'pending', 'expired'));
ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(255) NOT NULL;
ALTER TABLE invitations
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(255) NOT NULL;
ALTER TABLE users
    DROP COLUMN IF EXISTS password;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invitations
    DROP COLUMN IF EXISTS role_at_org;
ALTER TABLE invitations
    DROP COLUMN IF EXISTS department;
ALTER TABLE invitations
    DROP COLUMN IF EXISTS phone;
ALTER TABLE invitations
    DROP CONSTRAINT IF EXISTS invitations_status_check;
ALTER TABLE invitations
    DROP COLUMN IF EXISTS first_name;
ALTER TABLE invitations
    DROP COLUMN IF EXISTS last_name;
ALTER TABLE users
    ADD COLUMN password VARCHAR(255);
-- +goose StatementEnd
