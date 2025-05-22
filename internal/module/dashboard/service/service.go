package service

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/dashboard/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *dashboardService) GetDashboard(ctx context.Context) (*dto.GetDashboardResponse, error) {
	year := time.Now().Year()

	totalFaq, err := s.faqRepository.CountAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - Failed to get total FAQ")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalArticle, err := s.articleRepository.CountAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - Failed to get total article")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalCategory, err := s.articelCategoryRepository.CountAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - Failed to get total article category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalProductContent, err := s.productContentRepository.CountAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - Failed to get total product content")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	monthlyFaq := make([]dto.MonthlyCountFaq, 12)
	monthlyArticle := make([]dto.MonthlyCountArticle, 12)
	monthlyCategory := make([]dto.MonthlyCountCategory, 12)
	monthlyProductContent := make([]dto.MonthlyCountProductContent, 12)

	for m := 1; m <= 12; m++ {
		// start = 1st of month at 00:00, end = last of month 23:59
		start := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		monthName := start.Format("Jan")

		countFaq, err := s.faqRepository.CountByDate(ctx, start, end)
		if err != nil {
			log.Error().Err(err).Msg("service::GetDashboard - Failed to get count FAQ by date")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		countArt, err := s.articleRepository.CountByDate(ctx, start, end)
		if err != nil {
			log.Error().Err(err).Msg("service::GetDashboard - Failed to get count article by date")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		countCat, err := s.articelCategoryRepository.CountByDate(ctx, start, end)
		if err != nil {
			log.Error().Err(err).Msg("service::GetDashboard - Failed to get count article category by date")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		countPC, err := s.productContentRepository.CountByDate(ctx, start, end)
		if err != nil {
			log.Error().Err(err).Msg("service::GetDashboard - Failed to get count product content by date")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		monthlyFaq[m-1] = dto.MonthlyCountFaq{Month: monthName, CountFaq: countFaq}
		monthlyArticle[m-1] = dto.MonthlyCountArticle{Month: monthName, CountArticle: countArt}
		monthlyCategory[m-1] = dto.MonthlyCountCategory{Month: monthName, CountCategory: countCat}
		monthlyProductContent[m-1] = dto.MonthlyCountProductContent{Month: monthName, CountProductContent: countPC}
	}

	return &dto.GetDashboardResponse{
		FAQ: dto.FAQStats{
			TotalFaq:   totalFaq,
			MonthlyFaq: monthlyFaq,
		},
		Article: dto.ArticleStats{
			TotalArticle:   totalArticle,
			MonthlyArticle: monthlyArticle,
		},
		Category: dto.CategoryStats{
			TotalCategory:   totalCategory,
			MonthlyCategory: monthlyCategory,
		},
		ProductContent: dto.ProductContentStats{
			TotalProductContent:   totalProductContent,
			MonthlyProductContent: monthlyProductContent,
		},
	}, nil
}
