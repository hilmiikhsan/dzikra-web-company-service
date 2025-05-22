package route

import (
	article "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article/handler/rest"
	articleCategory "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/article_category/handler/rest"
	faq "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/faq/handler/rest"
	productContent "github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/module/product_content/handler/rest"
	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		// userAPI       = app.Group("/api/users")
		superadminAPI = app.Group("/api/superadmin")
		// publicAPI     = app.Group("/api")
	)

	productContent.NewProductContentHandler().ProductContentRoute(superadminAPI)
	faq.NewFAQHandler().FAQRoute(superadminAPI)
	articleCategory.NewArticleCategoryandler().ArticleCategoryRoute(superadminAPI)
	article.NewArticleHandler().ArticleRoute(superadminAPI)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
