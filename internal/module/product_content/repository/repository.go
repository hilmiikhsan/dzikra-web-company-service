package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductContentRepository = &productContentRepository{}

type productContentRepository struct {
	db *sqlx.DB
}

func NewProductContentRepository(db *sqlx.DB) *productContentRepository {
	return &productContentRepository{
		db: db,
	}
}

func (r *productContentRepository) InsertNewProductContent(ctx context.Context, tx *sqlx.Tx, data *entity.ProductContent) (*entity.ProductContent, error) {
	var res = new(entity.ProductContent)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductContent),
		data.ProductName,
		data.Images,
		data.ContentID,
		data.ContentEn,
		data.SellLink,
		data.WebLink,
	).Scan(
		&res.ID,
		&res.ProductName,
		&res.Images,
		&res.ContentID,
		&res.ContentEn,
		&res.SellLink,
		&res.WebLink,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::ContentID - Failed to insert new product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *productContentRepository) FindListProductContent(ctx context.Context, limit, offset int, search string) ([]dto.GetListProductContent, int, error) {
	var responses []entity.ProductContent

	args := []interface{}{
		search, search, search, search, search,
		limit, offset,
	}

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListProductContent), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductContent - error executing query")
		return nil, 0, err
	}

	var total int
	countArgs := []interface{}{
		search, search, search, search, search,
	}

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListProductContent), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListProductContent - error counting banner")
		return nil, 0, err
	}

	publicURL := config.Envs.MinioStorage.PublicURL

	productContent := make([]dto.GetListProductContent, 0, len(responses))
	for _, v := range responses {
		productContent = append(productContent, dto.GetListProductContent{
			ID:          v.ID,
			ProductName: v.ProductName,
			Images:      utils.FormatMediaPathURL(v.Images, publicURL),
			ContentID:   v.ContentID,
			ContentEn:   v.ContentEn,
			SellLink:    v.SellLink,
			WebLink:     v.WebLink,
			CreatedAt:   utils.FormatTime(v.CreatedAt),
		})
	}

	return productContent, total, nil
}

func (r *productContentRepository) UpdateProductContent(ctx context.Context, tx *sqlx.Tx, data *entity.ProductContent) (*entity.ProductContent, error) {
	var res = new(entity.ProductContent)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductContent),
		data.ProductName,
		data.Images,
		data.ContentID,
		data.ContentEn,
		data.SellLink,
		data.WebLink,
		data.ID,
	).Scan(
		&res.ID,
		&res.ProductName,
		&res.Images,
		&res.ContentID,
		&res.ContentEn,
		&res.SellLink,
		&res.WebLink,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateProductContent - address with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductContentNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateProductContent - Failed to update product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *productContentRepository) FindProductContentByID(ctx context.Context, id int) (*entity.ProductContent, error) {
	var res = new(entity.ProductContent)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindProductContentByID), id).Scan(
		&res.ID,
		&res.ProductName,
		&res.Images,
		&res.ContentID,
		&res.ContentEn,
		&res.SellLink,
		&res.WebLink,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::FindProductContentByID - address with id %d is not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductContentNotFound)
		}

		log.Error().Err(err).Msg("repository::FindProductContentByID - Failed to find product content by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *productContentRepository) SoftDeleteProductContentByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteProductContentByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteProductContentByID - Failed to soft delete product content")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteProductContentByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrProductContentNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteProductContentByID - Product Content not found")
		return errNotFound
	}

	return nil
}
