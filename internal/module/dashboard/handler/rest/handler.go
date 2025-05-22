package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *dashboardHandler) getDasbhboard(c *fiber.Ctx) error {
	res, err := h.service.GetDashboard(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("handler::getDasbhboard - Failed to get dashboard")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
