package rest

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productContentHandler) createProductContent(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductContentRequest)
		ctx = c.Context()
	)

	productName := c.FormValue("product_name")
	req.ProductName = productName
	contentID := c.FormValue("content_id")
	req.ContentID = contentID

	contentEn := c.FormValue("content_en")
	req.ContentEn = contentEn

	sellLink := c.FormValue("sell_link")
	req.SellLink = sellLink

	webLink := c.FormValue("web_link")
	req.WebLink = webLink

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductContent - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createAddress - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	mf, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to parse multipart form")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid multipart form data"))
	}

	fileHeaders := mf.File[constants.MultipartFormFile]
	switch len(fileHeaders) {
	case 0:
		log.Error().Msg("handler::createProductContent - No image file uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file is required"))
	case 1:
		log.Info().Msgf("handler::createProductContent - %s file is valid", fileHeaders[0].Filename)
	default:
		log.Error().Msgf("handler::createProductContent - too many files uploaded: %d", len(fileHeaders))
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
	}

	fh := fileHeaders[0]
	if fh.Size > constants.MaxFileSize {
		log.Error().Msg("handler::createProductContent - File size exceeds limit")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !constants.AllowedImageExtensions[ext] {
		log.Error().Msg("handler::createProductContent - Invalid file extension")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
	}

	file, err := fh.Open()
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to open file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}

	mimeType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		log.Error().Msg("handler::createProductContent - Uploaded file is not a valid image")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
	}

	objectName := "product_content_images/" + utils.GenerateBucketFileUUID() + ext
	uploadFile := dto.UploadFileRequest{
		ObjectName:     objectName,
		File:           fileBytes,
		FileHeaderSize: fh.Size,
		ContentType:    mimeType,
		Filename:       fh.Filename,
	}

	res, err := h.service.CreateProductContent(ctx, uploadFile, productName, contentID, contentEn, sellLink, webLink)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to create product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *productContentHandler) getListProductContent(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListProductContent(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProductContent - Failed to get list of product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productContentHandler) updateProductContent(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductContentRequest)
		ctx = c.Context()
	)

	productContentIDStr := c.Params("product_id")
	if productContentIDStr == "" {
		log.Warn().Msg("handler::updateProductContent - Product Content ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product Content ID is required"))
	}

	id, err := strconv.Atoi(productContentIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProductContent - Invalid product content ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product content ID"))
	}

	productName := c.FormValue("product_name")
	req.ProductName = productName
	contentID := c.FormValue("content_id")
	req.ContentID = contentID

	contentEn := c.FormValue("content_en")
	req.ContentEn = contentEn

	sellLink := c.FormValue("sell_link")
	req.SellLink = sellLink

	webLink := c.FormValue("web_link")
	req.WebLink = webLink

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductContent - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createAddress - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	var uploadFile *dto.UploadFileRequest
	if mf, err := c.MultipartForm(); err == nil {
		files := mf.File[constants.MultipartFormFile]
		if len(files) > 1 {
			log.Error().Msgf("handler::updateProductContent - too many files uploaded: %d", len(files))
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
		}
		if len(files) == 1 {
			fh := files[0]
			if fh.Size > constants.MaxFileSize {
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
			}
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !constants.AllowedImageExtensions[ext] {
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
			}
			f, err := fh.Open()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}
			mime := http.DetectContentType(data)
			if !strings.HasPrefix(mime, "image/") {
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
			}
			objectName := "product_content_images/" + utils.GenerateBucketFileUUID() + ext
			uploadFile = &dto.UploadFileRequest{
				ObjectName:     objectName,
				File:           data,
				FileHeaderSize: fh.Size,
				ContentType:    mime,
				Filename:       fh.Filename,
			}
		}
	}

	res, err := h.service.UpdateProductContent(ctx, uploadFile, productName, contentID, contentEn, sellLink, webLink, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductContent - Failed to create product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *productContentHandler) removeProductContent(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	productContentIDStr := c.Params("product_id")
	if productContentIDStr == "" {
		log.Warn().Msg("handler::removeProductContent - Product Content ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product Content ID is required"))
	}

	id, err := strconv.Atoi(productContentIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeProductContent - Invalid product content ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product content ID"))
	}

	if err := h.service.RemoveProductContent(ctx, id); err != nil {
		log.Error().Err(err).Msg("handler::removeProductContent - Failed to remove product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}

func (h *productContentHandler) getDetailProductContent(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	productContentIDStr := c.Params("product_id")
	if productContentIDStr == "" {
		log.Warn().Msg("handler::getDetailProductContent - Product Content ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product Content ID is required"))
	}

	id, err := strconv.Atoi(productContentIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailProductContent - Invalid product content ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product content ID"))
	}

	res, err := h.service.GetDetailProductContent(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailProductContent - Failed to get detail of product content")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
