package trails
import (
	"net/http"
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
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

	return router
}
