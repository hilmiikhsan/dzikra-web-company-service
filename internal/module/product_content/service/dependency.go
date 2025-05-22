package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/integration/storage/minio"
	productContentPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/ports"
	"github.com/jmoiron/sqlx"
)

var _ productContentPorts.ProductContentService = &productContentService{}

type productContentService struct {
	db                       *sqlx.DB
	productContentRepository productContentPorts.ProductContentRepository
	minioService             minio.MinioService
}

func NewProductContentService(
	db *sqlx.DB,
	productContentRepository productContentPorts.ProductContentRepository,
	minioService minio.MinioService,
) *productContentService {
	return &productContentService{
		db:                       db,
		productContentRepository: productContentRepository,
		minioService:             minioService,
	}
}
