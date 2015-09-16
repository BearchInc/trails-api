package controllers
import (
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"net/http"
	"github.com/bearchinc/trails-api/rest"
	"github.com/go-martini/martini"
	"time"
)

func TrailNextEvaluation(render render.Render, account *models.Account, db *appx.Datastore) {
	var trails = make([]*models.Trail, 0)
	if err := db.Query(models.Trails.ByNextEvaluation(account)).Results(&trails); err != nil {
		println("Next evaluation error: ", err.Error())

		count, err := db.Query(models.Trails.ByAccount(account)).Count()
		if count == 0 {
			println("Sleeping - I'm so lazy!!")
			time.Sleep(time.Second)
		}
		if err != nil {
			println("Next evaluation ---> Count error: ", err.Error())
		}

		TrailNextEvaluation(render, account, db)
		return
	}

	render.JSON(http.StatusOK, rest.FromTrails(trails))
}

func TrailLike(render render.Render, account *models.Account, db *appx.Datastore, params martini.Params) {
	println("Good job liking it")

	if err := models.Trails.Like(params["trail_id"], db); err != nil {
		println("The error: ", err.Error())
		render.JSON(http.StatusInternalServerError, err)
	}

	render.Status(http.StatusNoContent)
}

func TrailDislike(render render.Render, account *models.Account, db *appx.Datastore, params martini.Params) {
	println("You are evil")

	if err := models.Trails.Dislike(params["trail_id"], db); err != nil {
		println("The error: ", err.Error())
		render.JSON(http.StatusInternalServerError, err)
	}

	render.Status(http.StatusNoContent)

}