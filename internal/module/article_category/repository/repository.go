package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ArticleCategoryRepository = &articleCategoryRepository{}

type articleCategoryRepository struct {
	db *sqlx.DB
}

func NewArticleCategoryRepository(db *sqlx.DB) *articleCategoryRepository {
	return &articleCategoryRepository{
		db: db,
	}
}

func (r *articleCategoryRepository) InsertNewArticleCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ArticleCategory) (*entity.ArticleCategory, error) {
	var res = new(entity.ArticleCategory)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewArticleCategory),
		data.Name,
	).Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"article_categories_name_key": constants.ErrArticleCategoryAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewArticleCategory - Failed to insert new article category")
			return nil, handleErr
		}

		if adticleCategory, ok := val.(*entity.ArticleCategory); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewArticleCategory - Failed to insert new article category")
			return adticleCategory, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewArticleCategory - Failed to insert new article category")
		return nil, err
	}

	return res, nil
}

func (r *articleCategoryRepository) FindListArticleCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetListArticleCategory, int, error) {
	var responses []entity.ArticleCategory

	args := []interface{}{
		search, limit, offset,
	}

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListArticleCategory), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListArticleCategory - error executing query")
		return nil, 0, err
	}

	var total int
	countArgs := args[:1]

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListArticleCategory), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListArticleCategory - error counting article category")
		return nil, 0, err
	}

	articleCategories := make([]dto.GetListArticleCategory, 0, len(responses))
	for _, v := range responses {
		articleCategories = append(articleCategories, dto.GetListArticleCategory{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: utils.FormatTime(v.CreatedAt),
		})
	}

	return articleCategories, total, nil
}

func (r *articleCategoryRepository) FindArticleCategoryByID(ctx context.Context, id int) (*entity.ArticleCategory, error) {
	var res = new(entity.ArticleCategory)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindArticleCategoryByID), id).Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::FindArticleCategoryByID - address with id %d is not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrArticleCategoryNotFound)
		}

		log.Error().Err(err).Msg("repository::FindArticleCategoryByID - Failed to find article category by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *articleCategoryRepository) UpdateArticleCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ArticleCategory) (*entity.ArticleCategory, error) {
	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateArticleCategory),
		data.Name,
		data.ID,
	).Scan(
		&data.ID,
		&data.Name,
		&data.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateArticleCategory - article category with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrArticleCategoryNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateArticleCategory - Failed to update article category")
		return nil, err
	}

	return data, nil
}

func (r *articleCategoryRepository) SoftDeleteArticleCategoryByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteArticleCategoryByID), id)
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteArticleCategoryByID - Failed to soft delete article category")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteArticleCategoryByID - Failed to get rows affected")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if rowsAffected == 0 {
		errMessage := fmt.Errorf("repository::SoftDeleteArticleCategoryByID - article category with id %d is not found", id)
		log.Error().Err(err).Msg(errMessage.Error())
		return errors.New(constants.ErrArticleCategoryNotFound)
	}

	return nil
}

func (r *articleCategoryRepository) CountArticleCategoryByID(ctx context.Context, id int) (int, error) {
	var count int

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryCountArticleCategoryByID), id).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountArticleCategoryByID - Failed to count article category by id")
		return 0, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return count, nil
}

func (r *articleCategoryRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryCountAll)).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountAll - Failed to count all article categories")
		return 0, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return count, nil
}

func (r *articleCategoryRepository) CountByDate(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryCountByDate), start, end).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountByDate - Failed to count article categories by date")
		return 0, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return count, nil
}
