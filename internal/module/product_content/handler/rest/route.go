package rest

import "github.com/gofiber/fiber/v2"

func (h *productContentHandler) ProductContentRoute(superadminRouter, publicRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/product-content/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "product_content"), h.createProductContent)
	superadminRouter.Get("/product-content", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_content"), h.getListProductContent)
	superadminRouter.Patch("/product-content/update/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "product_content"), h.updateProductContent)
	superadminRouter.Delete("/product-content/remove/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "product_content"), h.removeProductContent)
	superadminRouter.Get("/product-content/:product_id", h.middleware.AuthBearer, h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_content"), h.getDetailProductContent)

	// public endpoint
	publicRouter.Get("/product-content", h.getListProductContent)
	publicRouter.Get("/product-content/:product_id", h.getDetailProductContent)
}
