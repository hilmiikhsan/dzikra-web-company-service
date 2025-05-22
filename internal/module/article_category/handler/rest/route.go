package rest

import "github.com/gofiber/fiber/v2"

func (h *articleCategoryHandler) ArticleCategoryRoute(superadminRouter, publicRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/category/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "category_article"), h.createArticleCategory)
	superadminRouter.Get("/category", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "category_article"), h.getListArticleCategory)
	superadminRouter.Get("/category/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "category_article"), h.getDetailArticleCategory)
	superadminRouter.Patch("/category/update/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "category_article"), h.updateArticleCategory)
	superadminRouter.Delete("/category/remove/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "category_article"), h.removeArticleCategory)

	// public endpoint
	publicRouter.Get("/category", h.getListArticleCategory)
	publicRouter.Get("/category/:category_id", h.getDetailArticleCategory)
}
