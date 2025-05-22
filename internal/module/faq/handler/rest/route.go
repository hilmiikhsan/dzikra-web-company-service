package rest

import "github.com/gofiber/fiber/v2"

func (h *faqHandler) FAQRoute(superadminRouter, publicRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/faq/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "faq"), h.createFAQ)
	superadminRouter.Get("/faq", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "faq"), h.getListFAQ)
	superadminRouter.Get("/faq/:faq_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "faq"), h.getDetailFAQ)
	superadminRouter.Patch("/faq/update/:faq_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "faq"), h.updateFAQ)
	superadminRouter.Delete("/faq/remove/:faq_id", h.middleware.AuthBearer, h.removeFAQ)

	// public endpoint
	publicRouter.Get("/faq", h.getListFAQ)
	publicRouter.Get("/faq/:faq_id", h.getDetailFAQ)
}
