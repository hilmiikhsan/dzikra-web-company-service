package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *articleCategoryHandler) createArticleCategory(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateArticleCategoryRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createArticleCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createArticleCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateArticleCategory(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createArticleCategory - Failed to create article category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *articleCategoryHandler) getListArticleCategory(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListArticleCategory(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListArticleCategory - Failed to get list article category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleCategoryHandler) getDetailArticleCategory(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("category_id")
	)

	if id == "" {
		log.Warn().Msg("handler::getDetailArticleCategory - Article Category ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article Category ID is required"))
	}

	articleCategoryID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailArticleCategory - Invalid article category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article category ID"))
	}

	res, err := h.service.GetDetailArticleCategory(ctx, articleCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailArticleCategory - Failed to get detail article category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleCategoryHandler) updateArticleCategory(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateArticleCategoryRequest)
		ctx = c.Context()
		id  = c.Params("category_id")
	)

	if id == "" {
		log.Warn().Msg("handler::updateArticleCategory - Article Category ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article Category ID is required"))
	}

	articleCategoryID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailArticleCategory - Invalid article category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article category ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateArticleCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateArticleCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateArticleCategory(ctx, req, articleCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateArticleCategory - Failed to update article category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *articleCategoryHandler) removeArticleCategory(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("category_id")
	)

	if id == "" {
		log.Warn().Msg("handler::removeArticleCategory - Article Category ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Article Category ID is required"))
	}

	articleCategoryID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailArticleCategory - Invalid article category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid article category ID"))
	}

	err = h.service.RemoveArticleCategory(ctx, articleCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeArticleCategory - Failed to remove article category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
