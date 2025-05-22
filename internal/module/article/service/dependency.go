package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/integration/storage/minio"
	articlePorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/ports"
	articleCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/ports"
	"github.com/jmoiron/sqlx"
)

var _ articlePorts.ArticleService = &articleService{}

type articleService struct {
	db                        *sqlx.DB
	articleRepository         articlePorts.ArticleRepository
	articleCategoryRepository articleCategoryPorts.ArticleCategoryRepository
	minioService              minio.MinioService
}

func NewArticleService(
	db *sqlx.DB,
	articleRepository articlePorts.ArticleRepository,
	articleCategoryRepository articleCategoryPorts.ArticleCategoryRepository,
	minioService minio.MinioService,
) *articleService {
	return &articleService{
		db:                        db,
		articleRepository:         articleRepository,
		articleCategoryRepository: articleCategoryRepository,
		minioService:              minioService,
	}
}
