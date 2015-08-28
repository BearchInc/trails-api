package models

import (
	"appengine/datastore"
	"github.com/bearchinc/hobo-api/api/utils"
	"github.com/drborges/appx"
)

type Account struct {
	appx.Model
	keySpec *appx.KeySpec

	Id string `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	AuthToken string `json:"auth_token"`
}

func (account *Account) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:      "Accounts",
		StringID:  account.Id,
		HasParent: false,
	}
}

func (model *Account) HasKey() bool {
	return model.Id != ""
}

var Accounts = struct {
	ByAuthToken func(authToken string) *datastore.Query
	New         func(accountId string) *Account
}{

	ByAuthToken: func(authToken string) *datastore.Query {
		return datastore.NewQuery(new(Account).KeySpec().Kind).Filter("AuthToken=", authToken).Limit(1)
	},

	New: func(accountId string) *Account {
		account := &Account{
			Id:        accountId,
			AuthToken: utils.GenerateToken(accountId),
		}

		return account
	},
}