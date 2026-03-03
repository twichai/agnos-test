-- +goose Up
-- +goose StatementBegin
ALTER TABLE hospitals
ADD COLUMN IF NOT EXISTS hospital_url TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS hospitals_hospital_url_key
ON hospitals (hospital_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS hospitals_hospital_url_key;

ALTER TABLE hospitals
DROP COLUMN IF EXISTS hospital_url;
-- +goose StatementEnd
