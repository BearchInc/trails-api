package middlewares
import (
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"github.com/go-martini/martini"
)

func ExternalServiceAuthorizationProvider(ds *appx.Datastore, martiniContext martini.Context, account *models.Account) {
	authorization := &models.ExternalServiceAuthorization{}
	authorization.SetParentKey(account.Key())

	if err := ds.Load(authorization); err != nil {
		panic(err)
	}

	martiniContext.Map(authorization)
}
