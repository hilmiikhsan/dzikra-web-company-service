package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/entity"
	"github.com/jmoiron/sqlx"
)

type ProductContentRepository interface {
	InsertNewProductContent(ctx context.Context, tx *sqlx.Tx, data *entity.ProductContent) (*entity.ProductContent, error)
	FindListProductContent(ctx context.Context, limit, offset int, search string) ([]dto.GetListProductContent, int, error)
	UpdateProductContent(ctx context.Context, tx *sqlx.Tx, data *entity.ProductContent) (*entity.ProductContent, error)
	FindProductContentByID(ctx context.Context, id int) (*entity.ProductContent, error)
	SoftDeleteProductContentByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type ProductContentService interface {
	CreateProductContent(ctx context.Context, payloadFile dto.UploadFileRequest, productName, contentID, contentEn, sellLink, webLink string) (*dto.CreateOrUpdateProductContentResponse, error)
	UpdateProductContent(ctx context.Context, payloadFile *dto.UploadFileRequest, productName, contentID, contentEn, sellLink, webLink string, id int) (*dto.CreateOrUpdateProductContentResponse, error)
	GetListProductContent(ctx context.Context, page, limit int, search string) (*dto.GetListProductContentResponse, error)
	RemoveProductContent(ctx context.Context, id int) error
	GetDetailProductContent(ctx context.Context, id int) (*dto.GetListProductContent, error)
}
