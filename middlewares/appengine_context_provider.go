package middlewares

import (
	"appengine"
	"github.com/go-martini/martini"
	"net/http"
)

// AppengineContextProvider provides an injectable and namespaced
// instance of appengine.Context
func AppengineContextProvider(c martini.Context, req *http.Request) {
	gae := appengine.NewContext(req)
	namespace := appengine.ModuleName(gae)
	context, err := appengine.Namespace(gae, namespace)
	if err != nil {
		panic(err)
	}
	c.Map(context)
}
