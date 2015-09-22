package controllers
import (
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"net/http"
	"github.com/bearchinc/trails-api/rest"
	"github.com/go-martini/martini"
	"time"
	"github.com/bearchinc/trails-api/middlewares"
	"appengine"
)

type RecursiveCount int

func TrailNextEvaluation(render render.Render, account *models.Account, log *middlewares.Logger, db *appx.Datastore) {
	log.Infof("Folks, I will start evaluating for account: %+v", account.FirstName)
	trails := trailNextEvaluation(render, account, log, db, 0)

	render.JSON(http.StatusOK, rest.FromTrails(trails))
}

func trailNextEvaluation(render render.Render, account *models.Account, log *middlewares.Logger, db *appx.Datastore, recursiveCount RecursiveCount) []*models.Trail {
	var trails = make([]*models.Trail, 0)
	if err := db.Query(models.Trails.ByNextEvaluation(account)).Results(&trails); err != nil {
		log.Errorf("Next evaluation error: ", err.Error())

		count, err := db.Query(models.Trails.ByAccount(account)).Count()
		if count == 0 {
			log.Infof("Sleeping - I'm so lazy!!")
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Errorf("Next evaluation ---> Count error: ", err.Error())
		}
		recursiveCount++

		if recursiveCount <= 10 {
			trails = trailNextEvaluation(render, account, log, db, recursiveCount)
			return trails
		}

		return trails
	}

	return trails
}

func TrailLike(render render.Render, context appengine.Context, account *models.Account, db *appx.Datastore, params martini.Params) {
	println("Good job liking it")

	if err := models.Trails.Like(params["trail_id"], db, context); err != nil {
		println("The error: ", err.Error())
		render.JSON(http.StatusInternalServerError, err)
	}

	render.Status(http.StatusNoContent)
}

func TrailDislike(render render.Render, context appengine.Context, account *models.Account, db *appx.Datastore, params martini.Params) {
	println("You are evil")

	if err := models.Trails.Dislike(params["trail_id"], db, context); err != nil {
		println("The error: ", err.Error())
		render.JSON(http.StatusInternalServerError, err)
	}

	render.Status(http.StatusNoContent)

}

func Stories(render render.Render) {
	println("fetching stories")

	stories := []models.Story{
		models.Story{
			ImagePath: "/Camera Uploads/2013-02-24 15.47.44.mov",
			Title: "Heck is love",
			Trails: []models.Trail{},
		},
		models.Story{
			ImagePath: "/Camera Uploads/2013-02-24 15.55.31.jpg",
			Title: "Heck dont hurt me",
			Trails: []models.Trail{},
		},
		models.Story{
			ImagePath: "/Camera Uploads/2013-02-24 18.32.33.jpg",
			Title: "Dont hurt me",
			Trails: []models.Trail{},
		},
		models.Story{
			ImagePath: "/Camera Uploads/2013-03-01 12.38.58.jpg",
			Title: "No more",
			Trails: []models.Trail{},
		},
		models.Story{
			ImagePath: "/Camera Uploads/2013-03-01 18.29.19.jpg",
			Title: "Lizard is cool",
			Trails: []models.Trail{},
		},
		models.Story{
			ImagePath: "/Camera Uploads/2013-03-02 11.05.37.jpg",
			Title: "Lizard dont hurt me",
			Trails: []models.Trail{},
		},
	}

	render.JSON(http.StatusOK, stories)
}