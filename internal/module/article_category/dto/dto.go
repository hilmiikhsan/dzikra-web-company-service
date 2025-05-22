package dto

type CreateOrUpdateArticleCategoryRequest struct {
	Name string `json:"name" validate:"required,max=100,min=2,xss_safe"`
}

type CreateOrUpdateArticleCategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type GetListArticleCategoryResponse struct {
	Category    []GetListArticleCategory `json:"category"`
	TotalPages  int                      `json:"total_page"`
	CurrentPage int                      `json:"current_page"`
	PageSize    int                      `json:"page_size"`
	TotalData   int                      `json:"total_data"`
}

type GetListArticleCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ArticleCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
