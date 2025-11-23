-- +goose Up
CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name varchar(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;

