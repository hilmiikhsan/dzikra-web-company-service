package ports

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/entity"
	"github.com/jmoiron/sqlx"
)

type ArticleCategoryRepository interface {
	InsertNewArticleCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ArticleCategory) (*entity.ArticleCategory, error)
	FindListArticleCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetListArticleCategory, int, error)
	FindArticleCategoryByID(ctx context.Context, id int) (*entity.ArticleCategory, error)
	UpdateArticleCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ArticleCategory) (*entity.ArticleCategory, error)
	SoftDeleteArticleCategoryByID(ctx context.Context, tx *sqlx.Tx, id int) error
	CountArticleCategoryByID(ctx context.Context, id int) (int, error)
	CountAll(ctx context.Context) (int64, error)
	CountByDate(ctx context.Context, start, end time.Time) (int64, error)
}

type ArticleCategoryService interface {
	CreateArticleCategory(ctx context.Context, req *dto.CreateOrUpdateArticleCategoryRequest) (*dto.CreateOrUpdateArticleCategoryResponse, error)
	GetListArticleCategory(ctx context.Context, page, limit int, search string) (*dto.GetListArticleCategoryResponse, error)
	GetDetailArticleCategory(ctx context.Context, id int) (*dto.GetListArticleCategory, error)
	UpdateArticleCategory(ctx context.Context, req *dto.CreateOrUpdateArticleCategoryRequest, id int) (*dto.CreateOrUpdateArticleCategoryResponse, error)
	RemoveArticleCategory(ctx context.Context, id int) error
}
