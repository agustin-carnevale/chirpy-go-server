-- +goose Up
CREATE TABLE chirps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  body TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Create an index on user_id for faster lookups
CREATE INDEX idx_chirps_user_id ON chirps(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_chirps_user_id;
DROP TABLE chirps;
