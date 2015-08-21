package middlewares

import (
	"appengine"
	"github.com/drborges/appx"
	"github.com/go-martini/martini"
)

func AppxProvider(c martini.Context, context appengine.Context) {
	c.Map(appx.NewDatastore(context))
}
