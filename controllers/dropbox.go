package controllers

import (
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/drborges/appx"
	"github.com/martini-contrib/render"
)

type RegisterDropboxForm struct {
	AccessToken		string		`json:"access_token" binding:"required"`
}

func RegisterDropbox(render render.Render, logger *middlewares.Logger, appx *appx.Datastore) {
	logger.Infof("You are in register dropbox")
}
