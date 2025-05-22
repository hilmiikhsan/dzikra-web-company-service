package dto

type CreateOrUpdateFAQRequest struct {
	QuestionID string `json:"question_id" validate:"required,min=2,xss_safe"`
	QuestionEn string `json:"question_en" validate:"required,min=2,xss_safe"`
	AnswerID   string `json:"answer_id" validate:"required,min=2,xss_safe"`
	AnswerEn   string `json:"answer_en" validate:"required,min=2,xss_safe"`
}

type CreateOrUpdateFAQResponse struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GetListFAQResponse struct {
	Faqs        []GetListFAQ `json:"faqs"`
	TotalPages  int          `json:"total_page"`
	CurrentPage int          `json:"current_page"`
	PageSize    int          `json:"page_size"`
	TotalData   int          `json:"total_data"`
}

type GetListFAQ struct {
	ID         int    `json:"id"`
	QuestionID string `json:"question_id"`
	QuestionEn string `json:"question_en"`
	AnswerID   string `json:"answer_id"`
	AnswerEn   string `json:"answer_en"`
	CreatedAt  string `json:"created_at"`
}
