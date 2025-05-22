package repository

const (
	queryInsertNewProductContent = `
		INSERT INTO product_contents
		(
			product_name,
			images, 
			content_id, 
			content_en, 
			sell_link, 
			web_link
		) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING 
			id, 
			product_name,
			images, 
			content_id, 
			content_en, 
			sell_link, 
			web_link, 
			created_at
	`

	queryFindListProductContent = `
		SELECT
			id,
			product_name,
			images, 
			content_id, 
			content_en, 
			sell_link, 
			web_link
		FROM product_contents
		WHERE
			deleted_at IS NULL
			AND (
                product_name ILIKE '%' || ? || '%'
                OR content_id   ILIKE '%' || ? || '%'
                OR content_en   ILIKE '%' || ? || '%'
                OR sell_link    ILIKE '%' || ? || '%'
                OR web_link     ILIKE '%' || ? || '%'
            )
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListProductContent = `
		SELECT COUNT(*)
		FROM product_contents
		WHERE 
			deleted_at IS NULL
			AND (
                product_name ILIKE '%' || ? || '%'
                OR content_id   ILIKE '%' || ? || '%'
                OR content_en   ILIKE '%' || ? || '%'
                OR sell_link    ILIKE '%' || ? || '%'
                OR web_link     ILIKE '%' || ? || '%'
            )
	`

	queryUpdateProductContent = `
		UPDATE product_contents
		SET
			product_name = ?,
			images = ?,
			content_id = ?,
			content_en = ?,
			sell_link = ?,
			web_link = ?,
			updated_at = NOW()
		WHERE 
			id = ?
			AND deleted_at IS NULL
		RETURNING
			id,
			product_name,
			images,
			content_id,
			content_en,
			sell_link,
			web_link,
			created_at 
	`

	queryFindProductContentByID = `
		SELECT
			id,
			product_name,
			images,
			content_id,
			content_en,
			sell_link,
			web_link,
			created_at
		FROM product_contents
		WHERE
			id = ?
			AND deleted_at IS NULL
	`

	querySoftDeleteProductContentByID = `
		UPDATE product_contents
		SET
			deleted_at = NOW()
		WHERE
			id = ?
			AND deleted_at IS NULL
	`

	queryCountAllProductContent = `
		SELECT COUNT(*) FROM product_contents WHERE deleted_at IS NULL
	`

	queryCountByDateProductContent = `
		SELECT COUNT(*) FROM product_contents 
		WHERE deleted_at IS NULL 
		AND created_at BETWEEN ? AND ?
	`
)
