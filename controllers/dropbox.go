package controllers

import (
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"net/http"
	"github.com/bearchinc/trails-api/services"
	"github.com/drborges/rivers/stream"
	"fmt"
)

type RegisterDropboxForm struct {
	AccessToken		string		`json:"access_token" binding:"required"`
	UserId			string		`json:"user_id" binding:"required"`
}

func RegisterDropbox(render render.Render, registerDropboxForm RegisterDropboxForm, account *models.Account, logger *middlewares.Logger, appx *appx.Datastore) {
	logger.Infof("You are in register dropbox")

	authorization := &models.ExternalServiceAuthorization{
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

func DropboxInit(req *http.Request, ds *appx.Datastore, account *models.Account, authorization *models.ExternalServiceAuthorization) {
	services.DropboxDelta(req, ds, authorization, newItem(ds))
}

func DropboxDelta(req *http.Request, ds *appx.Datastore, account *models.Account, authorization *models.ExternalServiceAuthorization) {
	services.DropboxDelta(req, ds, authorization, alreadyCategorized(ds))
}

func newItem(ds *appx.Datastore) stream.PredicateFn {
	return func(data stream.T) bool {
		return false
	}
}

func alreadyCategorized(ds *appx.Datastore) stream.PredicateFn {
	return func(data stream.T) bool {
		trail := data.(*models.Trail)

		err := ds.Load(trail)

		println(fmt.Sprint("Will I skip this value? : ", err == nil))
		return err == nil
	}
}
