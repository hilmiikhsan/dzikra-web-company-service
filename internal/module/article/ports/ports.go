package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/entity"
	"github.com/jmoiron/sqlx"
)

type ArticleRepository interface {
	InsertNewArticle(ctx context.Context, tx *sqlx.Tx, data *entity.Article) (*entity.Article, error)
	UpdateArticle(ctx context.Context, tx *sqlx.Tx, data *entity.Article) (*entity.Article, error)
	FindArticleByID(ctx context.Context, id int) (*entity.Article, error)
	FindListArticle(ctx context.Context, limit, offset int, search string) ([]dto.GetListArticle, int, error)
	SoftDeleteArticleByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type ArticleService interface {
	CreateArticle(ctx context.Context, payloadFile dto.UploadFileRequest, req *dto.CreateOrUpdateArticleRequest) (*dto.CreateOrUpdateArticleResponse, error)
	UpdateArticle(ctx context.Context, payloadFile *dto.UploadFileRequest, req *dto.CreateOrUpdateArticleRequest, id int) (*dto.CreateOrUpdateArticleResponse, error)
	GetListArticle(ctx context.Context, page, limit int, search string) (*dto.GetListArticleResponse, error)
	GetDetailArticle(ctx context.Context, id int) (*dto.GetListArticle, error)
	RemoveArticle(ctx context.Context, id int) error
}
