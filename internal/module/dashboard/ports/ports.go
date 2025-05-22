package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/dashboard/dto"
)

type DashboardService interface {
	GetDashboard(ctx context.Context) (*dto.GetDashboardResponse, error)
}
