package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

func (s *productContentService) CreateProductContent(ctx context.Context, payloadFile dto.UploadFileRequest, productName, contentID, contentEn, sellLink, webLink string) (*dto.CreateOrUpdateProductContentResponse, error) {
	// mapping file upload
	ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
	objectName := fmt.Sprintf("product_content_images/%s_%s", utils.GenerateBucketFileUUID(), ext)
	byteFile := utils.NewByteFile(payloadFile.File)

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", productName).Msg("service::CreateProductContent - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", productName).Msg("service::CreateProductContent - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ProductContent{
		ProductName: productName,
		Images:      objectName,
		ContentID:   contentID,
		ContentEn:   contentEn,
		SellLink:    sellLink,
		WebLink:     webLink,
	}

	res, err := s.productContentRepository.InsertNewProductContent(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Any("payload", payload).Msg("service::CreateProductContent - Failed to insert new product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload file to minio
	uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProductContent - Failed to upload file to minio")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("Uploaded image URL: %s", uploadedPath)

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateProductContent - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response
	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateProductContentResponse{
		ID:          res.ID,
		ProductName: res.ProductName,
		Images:      utils.FormatMediaPathURL(res.Images, publicURL),
		ContentID:   res.ContentID,
		ContentEn:   res.ContentEn,
		SellLink:    res.SellLink,
		WebLink:     res.WebLink,
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateProductContentResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productContentService) GetListProductContent(ctx context.Context, page, limit int, search string) (*dto.GetListProductContentResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list product content
	productContent, total, err := s.productContentRepository.FindListProductContent(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProductContent - error getting list product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if productContent is nil
	if productContent == nil {
		productContent = []dto.GetListProductContent{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListProductContentResponse{
		ProductContent: productContent,
		TotalPages:     totalPages,
		CurrentPage:    currentPage,
		PageSize:       perPage,
		TotalData:      total,
	}

	// return response
	return &response, nil
}

func (s *productContentService) UpdateProductContent(ctx context.Context, payloadFile *dto.UploadFileRequest, productName, contentID, contentEn, sellLink, webLink string, id int) (*dto.CreateOrUpdateProductContentResponse, error) {
	var (
		objectName string
		err        error
	)

	if payloadFile != nil {
		ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
		objectName = fmt.Sprintf("product_content_images/%s%s", utils.GenerateBucketFileUUID(), ext)
	}

	productContentResult, err := s.productContentRepository.FindProductContentByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductContentNotFound) {
			log.Error().Err(err).Msg("service::UpdateProductContent - Product content not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductContentNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateProductContent - Failed to find product content by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if productContentResult.Images != "" {
		// delete old banner image
		if err := s.minioService.DeleteFile(ctx, productContentResult.Images); err != nil {
			log.Error().Err(err).Msg("service::UpdateBanner - Failed to delete old product content image")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProductContent - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("service::UpdateProductContent - Failed to rollback transaction")
			}
		}
	}()

	payload := &entity.ProductContent{
		ID:          id,
		ProductName: productName,
		ContentID:   contentID,
		ContentEn:   contentEn,
		SellLink:    sellLink,
		WebLink:     webLink,
	}

	if objectName != "" {
		payload.Images = objectName
	}

	res, err := s.productContentRepository.UpdateProductContent(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductContentNotFound) {
			log.Error().Err(err).Msg("service::UpdateProductContent - Product content not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductContentNotFound))
		}

		log.Error().Err(err).Any("payload", payload).Msg("service::UpdateProductContent - Failed to update product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if payloadFile != nil {
		byteFile := utils.NewByteFile(payloadFile.File)
		if _, err = s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType); err != nil {
			log.Error().Err(err).Msg("service::UpdateProductContent - Failed to upload file to minio")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateProductContent - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	resp := &dto.CreateOrUpdateProductContentResponse{
		ID:          res.ID,
		ProductName: res.ProductName,
		Images:      utils.FormatMediaPathURL(res.Images, publicURL),
		ContentID:   res.ContentID,
		ContentEn:   res.ContentEn,
		SellLink:    res.SellLink,
		WebLink:     res.WebLink,
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}

	// Sanitasi output HTML
	policy := bluemonday.UGCPolicy()
	sanitized := utils.SanitizeCreateOrUpdateProductContentResponse(*resp, policy)
	return &sanitized, nil
}

func (s *productContentService) RemoveProductContent(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveProductContent - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveProductContent - Failed to rollback transaction")
			}
		}
	}()

	// soft delete product content
	if err := s.productContentRepository.SoftDeleteProductContentByID(ctx, tx, id); err != nil {
		if strings.Contains(err.Error(), constants.ErrProductContentNotFound) {
			log.Error().Err(err).Msg("service::RemoveProductContent - product content not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductContentNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveProductContent - Failed to soft delete product content")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveProductContent - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *productContentService) GetDetailProductContent(ctx context.Context, id int) (*dto.GetListProductContent, error) {
	productContent, err := s.productContentRepository.FindProductContentByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductContentNotFound) {
			log.Error().Err(err).Msg("service::GetDetailProductContent - product content not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductContentNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailProductContent - Failed to find product content by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.GetListProductContent{
		ID:          productContent.ID,
		ProductName: productContent.ProductName,
		Images:      utils.FormatMediaPathURL(productContent.Images, publicURL),
		ContentID:   productContent.ContentID,
		ContentEn:   productContent.ContentEn,
		SellLink:    productContent.SellLink,
		WebLink:     productContent.WebLink,
		CreatedAt:   utils.FormatTime(productContent.CreatedAt),
	}

	return response, nil
}
