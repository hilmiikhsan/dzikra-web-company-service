package dto

import (
	articleCategory "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/dto"
)

type CreateOrUpdateArticleRequest struct {
	Title       string `json:"title" validate:"required,max=100,min=2,xss_safe"`
	Description string `json:"desc" validate:"required,max=100,min=2,xss_safe"`
	Content     string `json:"content" validate:"required,min=1,xss_safe"`
	CategoryID  int    `json:"category_id" validate:"required,min=1"`
}

type CreateOrUpdateArticleResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	CategoryID  int    `json:"category_id"`
	Images      string `json:"images"`
	CreatedAt   string `json:"created_at"`
}

type UploadFileRequest struct {
	ObjectName     string `json:"object_name"`
	File           []byte `json:"-"`
	FileHeaderSize int64  `json:"-"`
	ContentType    string `json:"-"`
	Filename       string `json:"-"`
}

type GetListArticleResponse struct {
	Article     []GetListArticle `json:"article"`
	TotalPages  int              `json:"total_page"`
	CurrentPage int              `json:"current_page"`
	PageSize    int              `json:"page_size"`
	TotalData   int              `json:"total_data"`
}

type GetListArticle struct {
	ID          int                             `json:"id"`
	Title       string                          `json:"title"`
	Description string                          `json:"desc"`
	CategoryID  int                             `json:"category_id"`
	Images      string                          `json:"images"`
	Content     string                          `json:"content"`
	Category    articleCategory.ArticleCategory `json:"category"`
	CreatedAt   string                          `json:"created_at"`
}
