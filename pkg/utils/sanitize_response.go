package utils

import (
	productContent "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/dto"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeCreateOrUpdateProductContentResponse sanitizes the CreateOrUpdateBannerResponse by removing any potentially harmful content
func SanitizeCreateOrUpdateProductContentResponse(resp productContent.CreateOrUpdateProductContentResponse, policy *bluemonday.Policy) productContent.CreateOrUpdateProductContentResponse {
	resp.Images = policy.Sanitize(resp.Images)
	resp.ProductName = policy.Sanitize(resp.ProductName)
	resp.ContentID = policy.Sanitize(resp.ContentID)
	resp.ContentEn = policy.Sanitize(resp.ContentEn)
	resp.SellLink = policy.Sanitize(resp.SellLink)
	resp.WebLink = policy.Sanitize(resp.WebLink)

	return resp
}
