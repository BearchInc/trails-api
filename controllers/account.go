package controllers
import (
	"github.com/drborges/appx"
	"github.com/bearchinc/trails-api/models"
	"github.com/martini-contrib/render"
	"net/http"
)

type AccountUpdateForm struct {
	FirstName 		string `json:"first_name"`
	LastName  		string `json:"last_name"`
	Email     		string `json:"email"`
}

func UpdateAccount(render render.Render, accountUpdateForm AccountUpdateForm, account *models.Account, db *appx.Datastore) {
	account.FirstName = accountUpdateForm.FirstName
	account.LastName = accountUpdateForm.LastName
	account.Email = accountUpdateForm.Email

	if err := db.Save(account); err != nil {
		render.JSON(http.StatusInternalServerError, "Unable to register dropbox")
	}

	render.Status(http.StatusOK)
}