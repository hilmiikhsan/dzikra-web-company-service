package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/ports"
	faqRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/repository"
	faqService "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/service"
)

type faqHandler struct {
	service    ports.FAQService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewFAQHandler() *faqHandler {
	var handler = new(faqHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	faqRepository := faqRepository.NewFAQRepository(adapter.Adapters.DzikraPostgres)

	// product  service
	bannerService := faqService.NewFAQService(
		adapter.Adapters.DzikraPostgres,
		faqRepository,
	)

	// handler
	handler.service = bannerService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
