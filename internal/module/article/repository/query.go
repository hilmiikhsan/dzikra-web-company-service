package repository

const (
	queryInsertNewArticle = `
		INSERT INTO articles 
		(
			title, 
			image, 
			content, 
			description, 
			article_category_id
		) VALUES (?, ?, ?, ?, ?)
		RETURNING 
			id, 
			title, 
			image, 
			content, 
			description, 
			article_category_id, 
			created_at
	`

	queryUpdateArticle = `
		UPDATE articles
		SET
			title = ?,
			image = ?,
			content = ?,
			description = ?,
			article_category_id = ?
		WHERE id = ?
		RETURNING
			id,
			title,
			image,
			content,
			description,
			article_category_id,
			created_at
	`

	queryFindArticleByID = `
		SELECT
			a.id,
			a.title,
			a.image,
			a.content,
			a.description,
			a.article_category_id,
			c.name   AS article_category_name,
			a.created_at
		FROM articles a
		JOIN article_categories c
		ON a.article_category_id = c.id
		WHERE
			a.id = ?
			AND a.deleted_at IS NULL
	`

	queryFindListArticle = `
		SELECT
			a.id,
			a.title,
			a.image,
			a.content,
			a.description,
			a.article_category_id,
			c.name AS article_category_name,
			a.created_at
		FROM articles a
		JOIN article_categories c
		ON a.article_category_id = c.id
		WHERE
			a.deleted_at IS NULL
			AND (
				a.title       ILIKE '%' || ? || '%'
				OR a.content     ILIKE '%' || ? || '%'
				OR a.description ILIKE '%' || ? || '%'
				OR c.name        ILIKE '%' || ? || '%'
			)
		ORDER BY a.created_at DESC, a.id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListArticle = `
		SELECT COUNT(1)
		FROM articles a
		JOIN article_categories c
		ON a.article_category_id = c.id
		WHERE
			a.deleted_at IS NULL
			AND (
				a.title       ILIKE '%' || ? || '%'
				OR a.content     ILIKE '%' || ? || '%'
				OR a.description ILIKE '%' || ? || '%'
				OR c.name        ILIKE '%' || ? || '%'
			);
		`

	qyerySoftDeleteArticleByID = `
		UPDATE articles
		SET deleted_at = NOW()
		WHERE id = ?
		AND deleted_at IS NULL
	`

	queryCountAllArticle = `
		SELECT COUNT(*) FROM articles WHERE deleted_at IS NULL
	`

	queryCountByDateArticle = `
		SELECT COUNT(*) FROM articles 
		WHERE deleted_at IS NULL 
		AND created_at BETWEEN ? AND ?
	`
)
