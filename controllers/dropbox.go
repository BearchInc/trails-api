package controllers

import (
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"github.com/martini-contrib/render"
	"github.com/bearchinc/trails-api/models"
	"net/http"
	"github.com/bearchinc/trails-api/services"
	"github.com/drborges/rivers/stream"
	"appengine/datastore"
)

type RegisterDropboxForm struct {
	AccessToken string        `json:"access_token" binding:"required"`
	UserId      string        `json:"user_id" binding:"required"`
}

func RegisterDropbox(req *http.Request, render render.Render, registerDropboxForm RegisterDropboxForm, account *models.Account, logger *middlewares.Logger, ds *appx.Datastore) {
	logger.Infof("You are in register dropbox")

	authorization := &models.ExternalServiceAuthorization{
		AuthorizationType: models.DropBox,
		AccessToken: registerDropboxForm.AccessToken,
		UserId: models.DropBox.String() + "-" + registerDropboxForm.UserId,
	}

	authorization.SetParentKey(account.Key())

	err := ds.Load(authorization)
	if err != nil {
		println("I failed you becasue: %v", err.Error())
	}
	exists := err == nil

	if err := ds.Save(authorization); err != nil {
		logger.Errorf("Unable to register for dropbox %v", err)
		render.JSON(http.StatusInternalServerError, "Unable to register dropbox")
		return
	}

	if exists {
		DropboxDelta(req, ds, account, authorization)
	} else {
		DropboxInit(req, ds, account, authorization)
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
		if (err == datastore.ErrNoSuchEntity) {
			return false
		}

		return err == nil
	}
}
