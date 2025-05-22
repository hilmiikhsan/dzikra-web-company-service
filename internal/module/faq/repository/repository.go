package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.FAQRepository = &faqRepository{}

type faqRepository struct {
	db *sqlx.DB
}

func NewFAQRepository(db *sqlx.DB) *faqRepository {
	return &faqRepository{
		db: db,
	}
}

func (r *faqRepository) InsertNewFAQ(ctx context.Context, tx *sqlx.Tx, data *entity.FAQ) (*entity.FAQ, error) {
	var res = new(entity.FAQ)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewFAQ),
		data.Question,
		data.Answer,
	).Scan(
		&res.ID,
		&res.Question,
		&res.Answer,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewFAQ - Failed to insert new FAQ")
		return nil, err
	}

	return res, nil
}

func (r *faqRepository) FindListFAQ(ctx context.Context, limit, offset int, search string) ([]dto.GetListFAQ, int, error) {
	var responses []entity.FAQ

	args := []interface{}{
		search, search, search, search,
		limit, offset,
	}

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListFAQ), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListFAQ - error executing query")
		return nil, 0, err
	}

	var total int
	countArgs := args[:4]

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListFAQ), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListFAQ - error counting faq")
		return nil, 0, err
	}

	faqs := make([]dto.GetListFAQ, 0, len(responses))
	for _, v := range responses {
		qParts := strings.SplitN(v.Question, "|", 2)
		questionID, questionEn := "", ""
		if len(qParts) == 2 {
			questionID, questionEn = qParts[0], qParts[1]
		}

		aParts := strings.SplitN(v.Answer, "|", 2)
		answerID, answerEn := "", ""
		if len(aParts) == 2 {
			answerID, answerEn = aParts[0], aParts[1]
		}

		faqs = append(faqs, dto.GetListFAQ{
			ID:         v.ID,
			QuestionID: questionID,
			QuestionEn: questionEn,
			AnswerID:   answerID,
			AnswerEn:   answerEn,
			CreatedAt:  utils.FormatTime(v.CreatedAt),
		})
	}

	return faqs, total, nil
}

func (r *faqRepository) FindFAQByID(ctx context.Context, id int) (*entity.FAQ, error) {
	var res = new(entity.FAQ)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindFAQByID), id).Scan(
		&res.ID,
		&res.Question,
		&res.Answer,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::FindFAQByID - address with id %d is not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrFAQNotFound)
		}

		log.Error().Err(err).Msg("repository::FindFAQByID - Failed to find faq by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *faqRepository) UpdateFAQ(ctx context.Context, tx *sqlx.Tx, data *entity.FAQ) (*entity.FAQ, error) {
	var res = new(entity.FAQ)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateFAQ),
		data.Question,
		data.Answer,
		data.ID,
	).Scan(
		&res.ID,
		&res.Question,
		&res.Answer,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateFAQ - faq with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrFAQNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateFAQ - Failed to update faq")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *faqRepository) SoftDeleteFAQByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteFAQByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteFAQByID - Failed to soft delete faq")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteFAQByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrFAQNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteFAQByID - FAQ not found")
		return errNotFound
	}

	return nil
}
