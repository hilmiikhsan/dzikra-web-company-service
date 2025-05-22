package middleware

import (
	externalAuth "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/external/user"
)

type AuthMiddleware struct {
	externalAuth externalAuth.ExternalTokenValidation
}

func NewAuthMiddleware(externalAuth externalAuth.ExternalTokenValidation) *AuthMiddleware {
	return &AuthMiddleware{
		externalAuth: externalAuth,
	}
}
