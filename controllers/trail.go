package controllers
import (
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"net/http"
	"github.com/bearchinc/trails-api/rest"
)

func TrailNextEvaluation(render render.Render, account *models.Account, db *appx.Datastore) {
	var trails = make([]*models.Trail, 0)
	iter := db.Query(models.Trails.ByNextEvaluation(account)).PagesIterator()
	if err := iter.LoadNext(&trails); err != nil {
		println("The error: ", err.Error())
		render.JSON(http.StatusInternalServerError, err)
	}

	render.JSON(http.StatusOK, rest.FromTrails(trails, iter.Cursor()))
}

func TrailLike() {

}