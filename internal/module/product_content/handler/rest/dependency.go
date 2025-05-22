package rest

import (
	externalUser "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/integration/storage/minio"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/ports"
	productContentRepository "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/repository"
	productContentSerice "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/service"
)

type productContentHandler struct {
	service    ports.ProductContentService
	middleware middleware.AuthMiddleware
	validator  adapter.Validator
}

func NewProductContentHandler() *productContentHandler {
	var handler = new(productContentHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalAuth := &externalUser.External{}

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(externalAuth)

	// minio service
	minioService := minio.NewMinioService(adapter.Adapters.DzikraMinio, config.Envs.MinioStorage.Bucket)

	// repository
	productContentRepository := productContentRepository.NewProductContentRepository(adapter.Adapters.DzikraPostgres)

	// product  service
	bannerService := productContentSerice.NewProductContentService(
		adapter.Adapters.DzikraPostgres,
		productContentRepository,
		minioService,
	)

	// handler
	handler.service = bannerService
	handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}
