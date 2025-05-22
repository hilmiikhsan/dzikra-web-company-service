package repository

const (
	queryInsertNewFAQ = `
		INSERT INTO faqs 
		(
			question, 
			answer
		) VALUES (?, ?) RETURNING id, question, answer, created_at
	`

	queryFindListFAQ = `
		SELECT
			id,
			question,
			answer,
			created_at
		FROM faqs
		WHERE
			deleted_at IS NULL
			AND (
                split_part(question, '|', 1) ILIKE '%' || ? || '%'
                OR split_part(question, '|', 2) ILIKE '%' || ? || '%'
                OR split_part(answer,   '|', 1) ILIKE '%' || ? || '%'
                OR split_part(answer,   '|', 2) ILIKE '%' || ? || '%'
            )
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindListFAQ = `
		SELECT COUNT(*)
		FROM faqs
		WHERE 
			deleted_at IS NULL
			AND (
                split_part(question, '|', 1) ILIKE '%' || ? || '%'
                OR split_part(question, '|', 2) ILIKE '%' || ? || '%'
                OR split_part(answer,   '|', 1) ILIKE '%' || ? || '%'
                OR split_part(answer,   '|', 2) ILIKE '%' || ? || '%'
            )
	`

	queryFindFAQByID = `
		SELECT
			id,
			question,
			answer,
			created_at
		FROM faqs
		WHERE
			id = ?
			AND deleted_at IS NULL
	`

	queryUpdateFAQ = `
		UPDATE faqs
		SET
			question = ?,
			answer = ?
		WHERE
			id = ?
			AND deleted_at IS NULL
		RETURNING id, question, answer, created_at
	`

	querySoftDeleteFAQByID = `
		UPDATE faqs
		SET
			deleted_at = NOW()
		WHERE
			id = ?
			AND deleted_at IS NULL
	`
)
