package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *faqService) CreateFAQ(ctx context.Context, req *dto.CreateOrUpdateFAQRequest) (*dto.CreateOrUpdateFAQResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateFAQ - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateFAQ - Failed to rollback transaction")
			}
		}
	}()

	question := fmt.Sprintf("%s|%s", req.QuestionEn, req.QuestionID)
	answer := fmt.Sprintf("%s|%s", req.AnswerEn, req.AnswerID)

	payload := &entity.FAQ{
		Question: question,
		Answer:   answer,
	}

	res, err := s.faqRepository.InsertNewFAQ(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateFAQ - Failed to insert new FAQ")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateFAQ - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateFAQResponse{
		ID:       res.ID,
		Question: res.Question,
		Answer:   res.Answer,
	}, nil
}

func (s *faqService) GetListFAQ(ctx context.Context, page, limit int, search string) (*dto.GetListFAQResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list faq
	faqs, total, err := s.faqRepository.FindListFAQ(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListFAQ - error getting list faq")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if faqs is nil
	if faqs == nil {
		faqs = []dto.GetListFAQ{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListFAQResponse{
		Faqs:        faqs,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *faqService) GetDetailFAQ(ctx context.Context, id int) (*dto.GetListFAQ, error) {
	faq, err := s.faqRepository.FindFAQByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrFAQNotFound) {
			log.Error().Err(err).Msg("service::GetDetailFAQ - faq not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrFAQNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailFAQ - Failed to find faq by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	qParts := strings.SplitN(faq.Question, "|", 2)
	questionID, questionEn := "", ""
	if len(qParts) == 2 {
		questionID, questionEn = qParts[0], qParts[1]
	}

	aParts := strings.SplitN(faq.Answer, "|", 2)
	answerID, answerEn := "", ""
	if len(aParts) == 2 {
		answerID, answerEn = aParts[0], aParts[1]
	}

	response := &dto.GetListFAQ{
		ID:         faq.ID,
		QuestionID: questionID,
		QuestionEn: questionEn,
		AnswerID:   answerID,
		AnswerEn:   answerEn,
		CreatedAt:  utils.FormatTime(faq.CreatedAt),
	}

	return response, nil
}

func (s *faqService) UpdateFAQ(ctx context.Context, req *dto.CreateOrUpdateFAQRequest, id int) (*dto.CreateOrUpdateFAQResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateFAQ - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateFAQ - Failed to rollback transaction")
			}
		}
	}()

	question := fmt.Sprintf("%s|%s", req.QuestionEn, req.QuestionID)
	answer := fmt.Sprintf("%s|%s", req.AnswerEn, req.AnswerID)

	payload := &entity.FAQ{
		ID:       id,
		Question: question,
		Answer:   answer,
	}

	res, err := s.faqRepository.UpdateFAQ(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrFAQNotFound) {
			log.Error().Err(err).Msg("service::UpdateFAQ - faq not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrFAQNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateFAQ - Failed to update faq")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateFAQ - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateFAQResponse{
		ID:       res.ID,
		Question: res.Question,
		Answer:   res.Answer,
	}, nil
}

func (s *faqService) RemoveFAQ(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveFAQ - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveFAQ - Failed to rollback transaction")
			}
		}
	}()

	err = s.faqRepository.SoftDeleteFAQByID(ctx, tx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrFAQNotFound) {
			log.Error().Err(err).Msg("service::RemoveFAQ - faq not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrFAQNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveFAQ - Failed to soft delete faq")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveFAQ - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
