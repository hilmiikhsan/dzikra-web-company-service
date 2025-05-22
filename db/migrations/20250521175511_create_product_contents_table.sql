-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_contents (
    id SERIAL PRIMARY KEY,
    images TEXT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    content_id TEXT NOT NULL,
    content_en TEXT NOT NULL,
    sell_link VARCHAR(100) NOT NULL,
    web_link VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_contents_active  ON product_contents(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `product_contents`
CREATE TRIGGER set_updated_at_product_contents
BEFORE UPDATE ON product_contents
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_contents_active;
DROP TRIGGER IF EXISTS set_updated_at_product_contents ON product_contents;
DROP TABLE IF EXISTS product_contents;
-- +goose StatementEnd
