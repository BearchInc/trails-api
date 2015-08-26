package trails
import (
	"net/http"
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/bearchinc/trails-api/controllers"
	"github.com/martini-contrib/binding"
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

	router.Group("/account", func(r martini.Router){
		r.Get("/registerDropbox", controllers.RegisterDropbox)
	}, middlewares.AuthorizationAccountProvider)



	return router
}
