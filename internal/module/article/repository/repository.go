package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/ports"
	articleCategory "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/jmoiron/sqlx"
)

var _ ports.ArticleRepository = &articleRepository{}

type articleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) *articleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) InsertNewArticle(ctx context.Context, tx *sqlx.Tx, data *entity.Article) (*entity.Article, error) {
	var res = new(entity.Article)

	err := tx.QueryRowContext(ctx, tx.Rebind(queryInsertNewArticle),
		data.Title,
		data.Image,
		data.Content,
		data.Description,
		data.ArticleCategoryID,
	).Scan(
		&res.ID,
		&res.Title,
		&res.Image,
		&res.Content,
		&res.Description,
		&res.ArticleCategoryID,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewArticle - Failed to insert new article")
		return nil, err
	}

	return res, nil
}

func (R *articleRepository) UpdateArticle(ctx context.Context, tx *sqlx.Tx, data *entity.Article) (*entity.Article, error) {
	var res = new(entity.Article)

	err := tx.QueryRowContext(ctx, tx.Rebind(queryUpdateArticle),
		data.Title,
		data.Image,
		data.Content,
		data.Description,
		data.ArticleCategoryID,
		data.ID,
	).Scan(
		&res.ID,
		&res.Title,
		&res.Image,
		&res.Content,
		&res.Description,
		&res.ArticleCategoryID,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateArticle - article with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductContentNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateArticle - Failed to update article")
		return nil, err
	}

	return res, nil
}

func (r *articleRepository) FindArticleByID(ctx context.Context, id int) (*entity.Article, error) {
	var a entity.Article
	err := r.db.QueryRowxContext(ctx, r.db.Rebind(queryFindArticleByID), id).StructScan(&a)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("repository::FindArticleByID - article %d not found", id)
			return nil, errors.New(constants.ErrArticleNotFound)
		}

		log.Error().Err(err).Msg("repository::FindArticleByID - query failed")
		return nil, err
	}
	return &a, nil
}

func (r *articleRepository) FindListArticle(ctx context.Context, limit, offset int, search string) ([]dto.GetListArticle, int, error) {
	// pakai satu slice args untuk 4x search + limit+offset
	args := []interface{}{
		search, search, search, search,
		limit, offset,
	}

	var ents []entity.Article
	if err := r.db.SelectContext(ctx, &ents, r.db.Rebind(queryFindListArticle), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListArticle - error executing query")
		return nil, 0, err
	}

	// hitung total
	countArgs := []interface{}{search, search, search, search}
	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListArticle), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListArticle - error counting articles")
		return nil, 0, err
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	out := make([]dto.GetListArticle, 0, len(ents))
	for _, v := range ents {
		out = append(out, dto.GetListArticle{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Content:     v.Content,
			Images:      utils.FormatMediaPathURL(v.Image, publicURL),
			Category: articleCategory.ArticleCategory{
				ID:   v.ArticleCategoryID,
				Name: v.ArticleCategoryName,
			},
			CreatedAt: utils.FormatTime(v.CreatedAt),
		})
	}

	return out, total, nil
}

func (r *articleRepository) SoftDeleteArticleByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(qyerySoftDeleteArticleByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteArticleByID - Failed to soft delete article")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteArticleByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrArticleNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteArticleByID - article not found")
		return errNotFound
	}

	return nil
}

func (r *articleRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryCountAllArticle)).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountAll - error counting articles")
		return 0, err
	}

	return count, nil
}

func (r *articleRepository) CountByDate(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryCountByDateArticle), start, end).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountByDate - error counting articles by date")
		return 0, err
	}

	return count, nil
}
