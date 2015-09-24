package trails

import (
	"github.com/bearchinc/trails-api/controllers"
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
)

func Routes() http.Handler {
	// Middlewares are setup and run before each incoming request
	// the ones named like *Provider provide singleton instances
	// of injectable objects. For instance appx, logger, appengine context
	// can be injected in our routes handlers (a.k.a controllers)
	router := martini.Classic()

	router.Use(render.Renderer())
	router.Use(middlewares.AppengineContextProvider)
	router.Use(middlewares.LoggerProvider)
	router.Use(middlewares.AppxProvider)

	router.Post("/login", binding.Bind(controllers.LoginForm{}), controllers.Login)

	router.Group("/account", func(r martini.Router) {
		r.Post("/registerDropbox", binding.Bind(controllers.RegisterDropboxForm{}), controllers.RegisterDropbox)
		r.Post("/update", binding.Bind(controllers.AccountUpdateForm{}), controllers.UpdateAccount)

		r.Group("/dropbox", func(r martini.Router) {
			r.Post("/init", controllers.DropboxInit)
			r.Post("/delta", controllers.DropboxDelta)
		}, middlewares.ExternalServiceAuthorizationProvider)

		r.Group("/trails", func(r martini.Router) {
			r.Get("/next_evaluation", controllers.TrailNextEvaluation)
			r.Patch("/:trail_id/like", controllers.TrailLike)
			r.Patch("/:trail_id/dislike", controllers.TrailDislike)
			r.Get("/tags", controllers.Tags)
			r.Get("/tags/:tag_id", controllers.TagTrails)
		})
	}, middlewares.AuthorizationAccountProvider)

	return router
}
