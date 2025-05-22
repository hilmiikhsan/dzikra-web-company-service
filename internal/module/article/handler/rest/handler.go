package rest

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *articleHandler) createArticle(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateArticleRequest)
		ctx = c.Context()
	)

	title := c.FormValue("title")
	req.Title = title
	description := c.FormValue("desc")
	req.Description = description
	content := c.FormValue("content")
	req.Content = content
	categoryIDStr := c.FormValue("category_id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::createArticle - Invalid category_id")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid category_id"))
	}
	req.CategoryID = categoryID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createArticle - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createArticle - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	mf, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("handler::createArticle - Failed to parse multipart form")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid multipart form data"))
	}

	fileHeaders := mf.File[constants.MultipartFormFile]
	switch len(fileHeaders) {
	case 0:
		log.Error().Msg("handler::createArticle - No image file uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file is required"))
	case 1:
		log.Info().Msgf("handler::createArticle - %s file is valid", fileHeaders[0].Filename)
	default:
		log.Error().Msgf("handler::createArticle - too many files uploaded: %d", len(fileHeaders))
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
	}

	fh := fileHeaders[0]
	if fh.Size > constants.MaxFileSize {
		log.Error().Msg("handler::createArticle - File size exceeds limit")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !constants.AllowedImageExtensions[ext] {
		log.Error().Msg("handler::createArticle - Invalid file extension")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
	}

	file, err := fh.Open()
	if err != nil {
		log.Error().Err(err).Msg("handler::createArticle - Failed to open file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("handler::createArticle - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}

	mimeType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		log.Error().Msg("handler::createArticle - Uploaded file is not a valid image")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
	}

	objectName := "article_images/" + utils.GenerateBucketFileUUID() + ext
	uploadFile := dto.UploadFileRequest{
		ObjectName:     objectName,
		File:           fileBytes,
		FileHeaderSize: fh.Size,
		ContentType:    mimeType,
		Filename:       fh.Filename,
	}

	res, err := h.service.CreateArticle(ctx, uploadFile, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to create product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *articleHandler) updateArticle(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateArticleRequest)
		id  = c.Params("article_id")
		ctx = c.Context()
	)

	title := c.FormValue("title")
	req.Title = title
	description := c.FormValue("desc")
	req.Description = description
	content := c.FormValue("content")
	req.Content = content
	categoryIDStr := c.FormValue("category_id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateArticle - Invalid category_id")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid category_id"))
	}
	req.CategoryID = categoryID

	if id == "" {
		log.Warn().Msg("handler::updateArticle - Article ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article ID is required"))
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateArticle - Invalid article ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateArticle - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateArticle - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	var uploadFile *dto.UploadFileRequest
	if mf, err := c.MultipartForm(); err == nil {
		files := mf.File[constants.MultipartFormFile]
		if len(files) > 1 {
			log.Error().Msgf("handler::updateArticle - too many files uploaded: %d", len(files))
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
		}
		if len(files) == 1 {
			fh := files[0]
			if fh.Size > constants.MaxFileSize {
				log.Error().Msg("handler::updateArticle - File size exceeds limit")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
			}
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !constants.AllowedImageExtensions[ext] {
				log.Error().Msg("handler::updateArticle - Invalid file extension")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
			}
			f, err := fh.Open()
			if err != nil {
				log.Error().Err(err).Msg("handler::updateArticle - Failed to open file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				log.Error().Err(err).Msg("handler::updateArticle - Failed to read file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			mime := http.DetectContentType(data)
			if !strings.HasPrefix(mime, "image/") {
				log.Error().Msg("handler::updateArticle - Uploaded file is not a valid image")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
			}
			objectName := "article_images/" + utils.GenerateBucketFileUUID() + ext
			uploadFile = &dto.UploadFileRequest{
				ObjectName:     objectName,
				File:           data,
				FileHeaderSize: fh.Size,
				ContentType:    mime,
				Filename:       fh.Filename,
			}
		}
	}

	res, err := h.service.UpdateArticle(ctx, uploadFile, req, articleID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateArticle - Failed to update article")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleHandler) getListAticle(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListArticle(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListAticle - Failed to get list article")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleHandler) getDetailArticle(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("article_id")
	)

	if id == "" {
		log.Warn().Msg("handler::getDetailArticle - Article ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article ID is required"))
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailArticle - Invalid article ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article ID"))
	}

	res, err := h.service.GetDetailArticle(ctx, articleID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailArticle - Failed to get detail article")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleHandler) removeArticle(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("article_id")
	)

	if id == "" {
		log.Warn().Msg("handler::removeArticle - Article ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article ID is required"))
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeArticle - Invalid article ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article ID"))
	}

	err = h.service.RemoveArticle(ctx, articleID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeArticle - Failed to remove article")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
