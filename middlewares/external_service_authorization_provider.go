package middlewares
import (
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"github.com/go-martini/martini"
)

func ExternalServiceAuthorizationProvider(ds *appx.Datastore, martiniContext martini.Context) {
	authorization := &models.Authorization{}
	if err := ds.Load(authorization); err != nil {
		panic(err)
	}

	martiniContext.Map(&authorization)
}