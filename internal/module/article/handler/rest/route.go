package rest

import "github.com/gofiber/fiber/v2"

func (h *articleHandler) ArticleRoute(superadminRouter, publicRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/article/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "article"), h.createArticle)
	superadminRouter.Patch("/article/update/:article_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "article"), h.updateArticle)
	superadminRouter.Get("/article", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "article"), h.getListAticle)
	superadminRouter.Get("/article/:article_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "article"), h.getDetailArticle)
	superadminRouter.Delete("/article/remove/:article_id", h.middleware.AuthBearer, h.removeArticle)

	// public endpoint
	publicRouter.Get("/article", h.getListAticle)
	publicRouter.Get("/article/:article_id", h.getDetailArticle)
}
