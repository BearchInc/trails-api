package controllers

import (
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"net/http"
)

type RegisterDropboxForm struct {
	AccessToken		string		`json:"access_token" binding:"required"`
	UserId			string		`json:"user_id" binding:"required"`
}

func RegisterDropbox(render render.Render, registerDropboxForm RegisterDropboxForm, account *models.Account, logger *middlewares.Logger, appx *appx.Datastore) {
	logger.Infof("You are in register dropbox")

	authorization := &models.Authorization{
		AuthorizationType: models.DropBox,
		AccessToken: registerDropboxForm.AccessToken,
		UserId: models.DropBox.String() + "-" + registerDropboxForm.UserId,
	}


	authorization.SetParentKey(account.Key())

	if err := appx.Save(authorization); err != nil {
		logger.Errorf("Unable to register for dropbox %v", err)
		render.JSON(http.StatusInternalServerError, "Unable to register dropbox")
		return
	}

	render.Status(http.StatusOK)
}