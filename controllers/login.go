package controllers

import (
	"github.com/bearchinc/trails-api/middlewares"
	"github.com/bearchinc/trails-api/models"
	"github.com/drborges/appx"
	"github.com/martini-contrib/render"
	"net/http"
)

type LoginForm struct {
	Id string `json:"id" binding:"required"`
}

func Login(render render.Render, loginForm LoginForm, logger *middlewares.Logger, db *appx.Datastore, request *http.Request) {
	existingAccount := &models.Account{ Id: loginForm.Id }
	if err := db.Load(existingAccount); err == nil {
		render.JSON(http.StatusOK, existingAccount)
		return
	}

	account := models.Accounts.New(loginForm.Id)
	if err := db.Save(account); err != nil {
		logger.Errorf("Error while trying to login for account by id: %v", err.Error())
		render.JSON(http.StatusBadRequest, err.Error())
		return
	}

	render.JSON(http.StatusCreated, account)
}
