package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/integration/storage/minio"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/ports"
	articleRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/repository"
	articleService "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/service"
	articleCategoryRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/repository"
)

type articleHandler struct {
	service    ports.ArticleService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewArticleHandler() *articleHandler {
	var handler = new(articleHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// minio service
	minioService := minio.NewMinioService(adapter.Adapters.DzikraMinio, config.Envs.MinioStorage.Bucket)

	// repository
	articleRepository := articleRepository.NewArticleRepository(adapter.Adapters.DzikraPostgres)
	articleCategoryRepository := articleCategoryRepository.NewArticleCategoryRepository(adapter.Adapters.DzikraPostgres)

	// product  service
	bannerService := articleService.NewArticleService(
		adapter.Adapters.DzikraPostgres,
		articleRepository,
		articleCategoryRepository,
		minioService,
	)

	// handler
	handler.service = bannerService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
