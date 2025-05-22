package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *faqHandler) createFAQ(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateFAQRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createFAQ - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createFAQ - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateFAQ(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createFAQ - Failed to create faq")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *faqHandler) getListFAQ(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListFAQ(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListFAQ - Failed to get list faq")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *faqHandler) getDetailFAQ(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("faq_id")
	)

	if id == "" {
		log.Warn().Msg("handler::removeProductContent - Product Content ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product Content ID is required"))
	}

	faqID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailFAQ - Invalid faq ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid faq ID"))
	}

	res, err := h.service.GetDetailFAQ(ctx, faqID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailFAQ - Failed to get detail faq")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *faqHandler) updateFAQ(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("faq_id")
		req = new(dto.CreateOrUpdateFAQRequest)
	)

	if id == "" {
		log.Warn().Msg("handler::updateFAQ - FAQ ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("FAQ ID is required"))
	}

	faqID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateFAQ - Invalid faq ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid faq ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateFAQ - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateFAQ - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateFAQ(ctx, req, faqID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateFAQ - Failed to update faq")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *faqHandler) removeFAQ(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("faq_id")
	)

	if id == "" {
		log.Warn().Msg("handler::removeFAQ - FAQ ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("FAQ ID is required"))
	}

	faqID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeFAQ - Invalid faq ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid faq ID"))
	}

	err = h.service.RemoveFAQ(ctx, faqID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeFAQ - Failed to remove faq")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
