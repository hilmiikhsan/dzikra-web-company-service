-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS article_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_article_categories_name ON article_categories(name);
CREATE INDEX idx_article_categories_active  ON article_categories(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `article_categories`
CREATE TRIGGER set_updated_at_article_categories
BEFORE UPDATE ON article_categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_article_categories_name;
DROP INDEX IF EXISTS idx_article_categories_active;
DROP TRIGGER IF EXISTS set_updated_at_article_categories ON article_categories;
DROP TABLE IF EXISTS article_categories;
-- +goose StatementEnd
