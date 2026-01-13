-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role_at_mbl VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    password VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
	 phone VARCHAR(255) NOT NULL UNIQUE,
    date_of_birth DATE,
    state VARCHAR(100),
    department VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO users (first_name, last_name, role_at_mbl, role, email, phone, department)
VALUES ('Joshua', 'Ajagbe', 'Backend Engineer', 'admin', 'joshuaajagbe96@gmail.com', '+2348135860429', 'Backend');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

