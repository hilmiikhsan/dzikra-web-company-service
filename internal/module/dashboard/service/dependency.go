package service

import (
	articlePorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/ports"
	articleCategoryPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/ports"
	dashboardPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/dashboard/ports"
	faqPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/ports"
	productContentPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/ports"
	"github.com/jmoiron/sqlx"
)

var _ dashboardPorts.DashboardService = &dashboardService{}

type dashboardService struct {
	db                        *sqlx.DB
	faqRepository             faqPorts.FAQRepository
	articleRepository         articlePorts.ArticleRepository
	articelCategoryRepository articleCategoryPorts.ArticleCategoryRepository
	productContentRepository  productContentPorts.ProductContentRepository
}

func NewDashboardService(
	db *sqlx.DB,
	faqRepository faqPorts.FAQRepository,
	articleRepository articlePorts.ArticleRepository,
	articelCategoryRepository articleCategoryPorts.ArticleCategoryRepository,
	productContentRepository productContentPorts.ProductContentRepository,
) *dashboardService {
	return &dashboardService{
		db:                        db,
		faqRepository:             faqRepository,
		articleRepository:         articleRepository,
		articelCategoryRepository: articelCategoryRepository,
		productContentRepository:  productContentRepository,
	}
}
