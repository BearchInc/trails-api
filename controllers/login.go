package controllers
import (
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"net/http"
//	"github.com/martini-contrib/binding"
)

type LoginForm struct {
	AccountId string `json:"account_id" binding:"required"`
}

func Login(render render.Render, loginForm LoginForm, logger *middlewares.Logger, appx *appx.Datastore) {
	render.JSON(http.StatusOK, loginForm)
}

