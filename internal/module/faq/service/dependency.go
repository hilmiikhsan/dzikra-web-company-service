package service

import (
	faqPorts "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/ports"
	"github.com/jmoiron/sqlx"
)

var _ faqPorts.FAQService = &faqService{}

type faqService struct {
	db            *sqlx.DB
	faqRepository faqPorts.FAQRepository
}

func NewFAQService(
	db *sqlx.DB,
	faqRepository faqPorts.FAQRepository,
) *faqService {
	return &faqService{
		db:            db,
		faqRepository: faqRepository,
	}
}
