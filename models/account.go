package models
import (
	"github.com/drborges/appx"
	"appengine/datastore"
	"github.com/bearchinc/hobo-api/api/utils"
	"appengine"
)

type Account struct {
	appx.Model
	keySpec *appx.KeySpec

	Id				string	`json:"id"`

	FirstName		string	`json:"first_name"`
	LastName		string 	`json:"last_name"`
	Email			string	`json:"email"`

	AuthToken		string	`json:"auth_token"`
}

func (account *Account) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind: 		"Accounts",
		StringID: account.Id,
		HasParent: false,
	}
}

func (model *Account) HasKey() bool {
	return model.Id != ""
}

var Accounts = struct {
	ByAuthToken			func(authToken string) *datastore.Query
	New					func(accountId string) *Account
	NewWithAuthToken	func(accountId string, context appengine.Context) *Account
} {

	ByAuthToken: func(authToken string) *datastore.Query {
		return datastore.NewQuery(new(Account).KeySpec().Kind).Filter("AuthToken=", authToken).Limit(1)
	},

	New: func(accountId string) *Account {
		return &Account { Id: accountId }
	},

	NewWithAuthToken: func(accountId string, context appengine.Context) *Account {
		account := &Account{
			Id: accountId,
			AuthToken: utils.GenerateToken(accountId),
		}

		if err := appx.NewKeyResolver(context).Resolve(account); err != nil {
			panic("You fucked up " + err.Error())
		}


		account.EncodedKey()

		return account
	},

}