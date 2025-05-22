package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/ports"
	articleCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/repository"
	articleCategoryService "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/service"
)

type articleCategoryHandler struct {
	service    ports.ArticleCategoryService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewArticleCategoryandler() *articleCategoryHandler {
	var handler = new(articleCategoryHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// repository
	articleCategoryRepository := articleCategoryRepository.NewArticleCategoryRepository(adapter.Adapters.DzikraPostgres)

	// article category service
	articleCategoryService := articleCategoryService.NewArticleCategoryService(
		adapter.Adapters.DzikraPostgres,
		articleCategoryRepository,
	)

	// handler
	handler.service = articleCategoryService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
