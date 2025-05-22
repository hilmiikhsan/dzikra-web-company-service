package rest

import "github.com/gofiber/fiber/v2"

func (h *dashboardHandler) DashboardRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Get("/dashboard", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "dashboard"), h.getDasbhboard)
}
