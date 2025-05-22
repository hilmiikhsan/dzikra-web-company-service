package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/entity"
	"github.com/jmoiron/sqlx"
)

type FAQRepository interface {
	InsertNewFAQ(ctx context.Context, tx *sqlx.Tx, data *entity.FAQ) (*entity.FAQ, error)
	FindListFAQ(ctx context.Context, limit, offset int, search string) ([]dto.GetListFAQ, int, error)
	FindFAQByID(ctx context.Context, id int) (*entity.FAQ, error)
	UpdateFAQ(ctx context.Context, tx *sqlx.Tx, data *entity.FAQ) (*entity.FAQ, error)
	SoftDeleteFAQByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type FAQService interface {
	CreateFAQ(ctx context.Context, req *dto.CreateOrUpdateFAQRequest) (*dto.CreateOrUpdateFAQResponse, error)
	GetListFAQ(ctx context.Context, page, limit int, search string) (*dto.GetListFAQResponse, error)
	GetDetailFAQ(ctx context.Context, id int) (*dto.GetListFAQ, error)
	UpdateFAQ(ctx context.Context, req *dto.CreateOrUpdateFAQRequest, id int) (*dto.CreateOrUpdateFAQResponse, error)
	RemoveFAQ(ctx context.Context, id int) error
}
