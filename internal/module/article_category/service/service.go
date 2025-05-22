package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *articleCategoryService) CreateArticleCategory(ctx context.Context, req *dto.CreateOrUpdateArticleCategoryRequest) (*dto.CreateOrUpdateArticleCategoryResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticleCategory - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateArticleCategory - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ArticleCategory{
		Name: req.Name,
	}

	res, err := s.articleCategoryRepository.InsertNewArticleCategory(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleCategoryAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateArticleCategory - Article category already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrArticleCategoryAlreadyRegistered))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticleCategory - Failed to insert new article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateArticleCategory - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateArticleCategoryResponse{
		ID:        res.ID,
		Name:      res.Name,
		CreatedAt: utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *articleCategoryService) GetListArticleCategory(ctx context.Context, page, limit int, search string) (*dto.GetListArticleCategoryResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list article category
	articleCategories, total, err := s.articleCategoryRepository.FindListArticleCategory(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListArticleCategory - error getting list article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if articleCategories is nil
	if articleCategories == nil {
		articleCategories = []dto.GetListArticleCategory{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListArticleCategoryResponse{
		Category:    articleCategories,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *articleCategoryService) GetDetailArticleCategory(ctx context.Context, id int) (*dto.GetListArticleCategory, error) {
	articleCategory, err := s.articleCategoryRepository.FindArticleCategoryByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleCategoryNotFound) {
			log.Error().Err(err).Msg("service::GetDetailArticleCategory - article category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailArticleCategory - error getting article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetListArticleCategory{
		ID:        articleCategory.ID,
		Name:      articleCategory.Name,
		CreatedAt: utils.FormatTime(articleCategory.CreatedAt),
	}, nil
}

func (s *articleCategoryService) UpdateArticleCategory(ctx context.Context, req *dto.CreateOrUpdateArticleCategoryRequest, id int) (*dto.CreateOrUpdateArticleCategoryResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateArticleCategory - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateArticleCategory - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ArticleCategory{
		ID:   id,
		Name: req.Name,
	}

	res, err := s.articleCategoryRepository.UpdateArticleCategory(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleCategoryNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateArticleCategory - article category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleCategoryNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateArticleCategory - Failed to update article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateArticleCategory - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateArticleCategoryResponse{
		ID:        res.ID,
		Name:      res.Name,
		CreatedAt: utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *articleCategoryService) RemoveArticleCategory(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveArticleCategory - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveArticleCategory - Failed to rollback transaction")
			}
		}
	}()

	err = s.articleCategoryRepository.SoftDeleteArticleCategoryByID(ctx, tx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleCategoryNotFound) {
			log.Error().Err(err).Msg("service::RemoveArticleCategory - article category not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveArticleCategory - Failed to soft delete article category")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveArticleCategory - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
