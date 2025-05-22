package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/entity"
	articleCategory "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

func (s *articleService) CreateArticle(ctx context.Context, payloadFile dto.UploadFileRequest, req *dto.CreateOrUpdateArticleRequest) (*dto.CreateOrUpdateArticleResponse, error) {
	countArticleCategory, err := s.articleCategoryRepository.CountArticleCategoryByID(ctx, req.CategoryID)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticle - Failed to count article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if countArticleCategory == 0 {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticle - Article category not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleCategoryNotFound))
	}

	// mapping file upload
	ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
	objectName := fmt.Sprintf("article_images/%s_%s", utils.GenerateBucketFileUUID(), ext)
	byteFile := utils.NewByteFile(payloadFile.File)

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticle - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateArticle - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.Article{
		Title:             req.Title,
		Image:             objectName,
		Content:           req.Content,
		Description:       req.Description,
		ArticleCategoryID: req.CategoryID,
	}

	res, err := s.articleRepository.InsertNewArticle(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateArticle - Failed to insert new article")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload file to minio
	uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateArticle - Failed to upload file to minio")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("Uploaded image URL: %s", uploadedPath)

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateArticle - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response
	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateArticleResponse{
		ID:          res.ID,
		Title:       res.Title,
		Description: res.Description,
		CategoryID:  res.ArticleCategoryID,
		Images:      utils.FormatMediaPathURL(res.Image, publicURL),
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateArticleResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *articleService) UpdateArticle(ctx context.Context, payloadFile *dto.UploadFileRequest, req *dto.CreateOrUpdateArticleRequest, id int) (*dto.CreateOrUpdateArticleResponse, error) {
	var (
		err error
	)

	countCat, err := s.articleCategoryRepository.CountArticleCategoryByID(ctx, req.CategoryID)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateArticle - Failed to count article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if countCat == 0 {
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleCategoryNotFound))
	}

	existing, err := s.articleRepository.FindArticleByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleNotFound) {
			log.Error().Err(err).Msg("service::UpdateArticle - article not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleNotFound))
		}
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var objectName string
	if payloadFile != nil {
		ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
		objectName = fmt.Sprintf("article_images/%s%s", utils.GenerateBucketFileUUID(), ext)

		if existing.Image != "" {
			if err := s.minioService.DeleteFile(ctx, existing.Image); err != nil {
				log.Error().Err(err).Msg("service::UpdateArticle - Failed to delete old article image")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}
	}

	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateArticle - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	payload := &entity.Article{
		ID:                id,
		Title:             req.Title,
		Description:       req.Description,
		Content:           req.Content,
		ArticleCategoryID: req.CategoryID,
		Image:             existing.Image,
	}

	if objectName != "" {
		payload.Image = objectName
	}

	res, err := s.articleRepository.UpdateArticle(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleNotFound) {
			log.Error().Err(err).Msg("service::UpdateArticle - article not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateArticle - Failed to update article")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if payloadFile != nil {
		byteFile := utils.NewByteFile(payloadFile.File)
		if _, err = s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType); err != nil {
			log.Error().Err(err).Msg("service::UpdateArticle - Failed to upload file to minio")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateArticle - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateArticleResponse{
		ID:          res.ID,
		Title:       res.Title,
		Description: res.Description,
		CategoryID:  res.ArticleCategoryID,
		Images:      utils.FormatMediaPathURL(res.Image, publicURL),
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateArticleResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *articleService) GetListArticle(ctx context.Context, page, limit int, search string) (*dto.GetListArticleResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list article
	articles, total, err := s.articleRepository.FindListArticle(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListArticle - error getting list article")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if articles is nil
	if articles == nil {
		articles = []dto.GetListArticle{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListArticleResponse{
		Article:     articles,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *articleService) GetDetailArticle(ctx context.Context, id int) (*dto.GetListArticle, error) {
	article, err := s.articleRepository.FindArticleByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleNotFound) {
			log.Error().Err(err).Msg("service::GetDetailArticle - article not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleNotFound))
		}
		log.Error().Err(err).Msg("service::GetDetailArticle - error getting article")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.GetListArticle{
		ID:          article.ID,
		Title:       article.Title,
		Description: article.Description,
		Content:     article.Content,
		Images:      utils.FormatMediaPathURL(article.Image, publicURL),
		CategoryID:  article.ArticleCategoryID,
		Category: articleCategory.ArticleCategory{
			ID:   article.ArticleCategoryID,
			Name: article.ArticleCategoryName,
		},
	}

	return response, nil
}

func (s *articleService) RemoveArticle(ctx context.Context, id int) error {
	article, err := s.articleRepository.FindArticleByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrArticleNotFound) {
			log.Error().Err(err).Msg("service::RemoveArticle - article not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrArticleNotFound))
		}
		log.Error().Err(err).Msg("service::RemoveArticle - error getting article")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveArticle - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveArticle - Failed to rollback transaction")
			}
		}
	}()

	if err := s.articleRepository.SoftDeleteArticleByID(ctx, tx, id); err != nil {
		log.Error().Err(err).Msg("service::RemoveArticle - error removing article")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if article.Image != "" {
		if err := s.minioService.DeleteFile(ctx, article.Image); err != nil {
			log.Error().Err(err).Msg("service::RemoveArticle - Failed to delete article image")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveArticle - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
