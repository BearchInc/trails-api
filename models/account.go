package models
import (
	"github.com/drborges/appx"
	"appengine/datastore"
)

type Account struct {
	appx.Model

	FirstName		string	`json:"first_name"`
	LastName		string 	`json:"last_name"`
	Email			string	`json:"email"`

	AuthToken		string	`json:"auth_token"`
}

func (account *Account) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind: 		"Accounts",
		Incomplete: false,
	}
}

var Accounts = struct {
	ByAuthToken		func(authToken string) *datastore.Query
} {

	ByAuthToken: func(authToken string) *datastore.Query {
		return datastore.NewQuery(new(Account).KeySpec().Kind).Filter("AuthToken=", authToken).Limit(1)
	},

}