package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/middleware"
	articleRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/repository"
	articleCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/dashboard/ports"
	dashboardService "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/dashboard/service"
	faqRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/repository"
	productContentRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/repository"
)

type dashboardHandler struct {
	service    ports.DashboardService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewDashboardHandler() *dashboardHandler {
	var handler = new(dashboardHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	faqRepository := faqRepository.NewFAQRepository(adapter.Adapters.DzikraPostgres)
	articleRepository := articleRepository.NewArticleRepository(adapter.Adapters.DzikraPostgres)
	articleCategoryRepository := articleCategoryRepository.NewArticleCategoryRepository(adapter.Adapters.DzikraPostgres)
	productContentRepository := productContentRepository.NewProductContentRepository(adapter.Adapters.DzikraPostgres)

	// dashboard service
	dashboardService := dashboardService.NewDashboardService(
		adapter.Adapters.DzikraPostgres,
		faqRepository,
		articleRepository,
		articleCategoryRepository,
		productContentRepository,
	)

	// handler
	handler.service = dashboardService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
