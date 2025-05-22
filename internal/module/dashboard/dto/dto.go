package dto

type GetDashboardResponse struct {
	FAQ            FAQStats            `json:"faq"`
	Article        ArticleStats        `json:"article"`
	Category       CategoryStats       `json:"category"`
	ProductContent ProductContentStats `json:"productContent"`
}

type FAQStats struct {
	TotalFaq   int64             `json:"total_faq"`
	MonthlyFaq []MonthlyCountFaq `json:"monthly_faq"`
}

type MonthlyCountFaq struct {
	Month    string `json:"month"`
	CountFaq int64  `json:"countFaq"`
}

type ArticleStats struct {
	TotalArticle   int64                 `json:"total_article"`
	MonthlyArticle []MonthlyCountArticle `json:"monthly_article"`
}

type MonthlyCountArticle struct {
	Month        string `json:"month"`
	CountArticle int64  `json:"countArticle"`
}

type CategoryStats struct {
	TotalCategory   int64                  `json:"total_category"`
	MonthlyCategory []MonthlyCountCategory `json:"monthly_category"`
}

type MonthlyCountCategory struct {
	Month         string `json:"month"`
	CountCategory int64  `json:"countCategory"`
}

type ProductContentStats struct {
	TotalProductContent   int64                        `json:"total_product_content"`
	MonthlyProductContent []MonthlyCountProductContent `json:"monthly_product_content"`
}

type MonthlyCountProductContent struct {
	Month               string `json:"month"`
	CountProductContent int64  `json:"countProductContent"`
}
