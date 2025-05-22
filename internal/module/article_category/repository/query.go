package repository

const (
	queryInsertNewArticleCategory = `
		INSERT INTO article_categories (name) VALUES (?) RETURNING id, name, created_at
	`

	queryFindListArticleCategory = `
		SELECT id, name, created_at
		FROM article_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListArticleCategory = `
		SELECT COUNT(*) FROM article_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`

	queryFindArticleCategoryByID = `
		SELECT id, name, created_at
		FROM article_categories
		WHERE id = ? AND deleted_at IS NULL
	`

	queryUpdateArticleCategory = `
		UPDATE article_categories
		SET name = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, created_at
	`

	querySoftDeleteArticleCategoryByID = `
		UPDATE article_categories
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryCountArticleCategoryByID = `
		SELECT COUNT(id) FROM article_categories
		WHERE id = ? AND deleted_at IS NULL
	`
)
