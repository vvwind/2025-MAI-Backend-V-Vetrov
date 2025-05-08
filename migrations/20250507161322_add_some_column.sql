-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN category TEXT DEFAULT 'no_category';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN IF EXISTS category;
-- +goose StatementEnd
