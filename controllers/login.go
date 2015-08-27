package controllers
import (
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"net/http"
	"github.com/bearchinc/trails-api/models"
	"appengine"
)

type LoginForm struct {
	AccountId string `json:"account_id" binding:"required"`
}

func Login(render render.Render, loginForm LoginForm, logger *middlewares.Logger, db *appx.Datastore, request *http.Request) {

	context := appengine.NewContext(request)

	existingAccount := models.Accounts.New(loginForm.AccountId)
	err := db.Load(existingAccount)

	if err == nil {
		render.JSON(http.StatusOK, existingAccount)
		return
	}

	account := models.Accounts.NewWithAuthToken(loginForm.AccountId, context)
	if err = db.Save(account); err == nil {
		render.JSON(http.StatusCreated, account)
		return
	}

	logger.Errorf("Error while trying to login for account by id: %v", err.Error())
	render.JSON(http.StatusBadRequest, err.Error())
}

