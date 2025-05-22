package service

import (
	articleCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/ports"
	"github.com/jmoiron/sqlx"
)

var _ articleCategoryPorts.ArticleCategoryService = &articleCategoryService{}

type articleCategoryService struct {
	db                        *sqlx.DB
	articleCategoryRepository articleCategoryPorts.ArticleCategoryRepository
}

func NewArticleCategoryService(
	db *sqlx.DB,
	articleCategoryRepository articleCategoryPorts.ArticleCategoryRepository,
) *articleCategoryService {
	return &articleCategoryService{
		db:                        db,
		articleCategoryRepository: articleCategoryRepository,
	}
}
