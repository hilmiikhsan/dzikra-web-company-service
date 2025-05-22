package dto

type CreateOrUpdateProductContentResponse struct {
	ID          int    `json:"id"`
	ProductName string `json:"product_name"`
	Images      string `json:"images"`
	ContentID   string `json:"content_id"`
	ContentEn   string `json:"content_en"`
	SellLink    string `json:"sell_link"`
	WebLink     string `json:"web_link"`
	CreatedAt   string `json:"created_at"`
}

type UploadFileRequest struct {
	ObjectName     string `json:"object_name"`
	File           []byte `json:"-"`
	FileHeaderSize int64  `json:"-"`
	ContentType    string `json:"-"`
	Filename       string `json:"-"`
}

type CreateOrUpdateProductContentRequest struct {
	ProductName string `json:"product_name" validate:"required,max=100,min=2,xss_safe"`
	ContentID   string `json:"content_id" validate:"required,max=100,min=2,xss_safe"`
	ContentEn   string `json:"content_en" validate:"required,max=100,min=2,xss_safe"`
	SellLink    string `json:"sell_link" validate:"required,url,xss_safe"`
	WebLink     string `json:"web_link" validate:"required,url,xss_safe"`
}

type GetListProductContentResponse struct {
	ProductContent []GetListProductContent `json:"product_content"`
	TotalPages     int                     `json:"total_page"`
	CurrentPage    int                     `json:"current_page"`
	PageSize       int                     `json:"page_size"`
	TotalData      int                     `json:"total_data"`
}

type GetListProductContent struct {
	ID          int    `json:"id"`
	ProductName string `json:"product_name"`
	Images      string `json:"images"`
	ContentID   string `json:"content_id"`
	ContentEn   string `json:"content_en"`
	SellLink    string `json:"sell_link"`
	WebLink     string `json:"web_link"`
	CreatedAt   string `json:"created_at"`
}
