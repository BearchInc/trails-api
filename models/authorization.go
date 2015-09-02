package models

import (
	"github.com/drborges/appx"
	"strconv"
	"appengine/datastore"
)

type AuthorizationType int

func (typ AuthorizationType) String() string {
	return strconv.Itoa(int(typ))
}

const (
	DropBox AuthorizationType = 1
)

type Authorization struct {
	appx.Model

	AuthorizationType AuthorizationType `json:"authorization_type"`
	AccessToken       string            `json:"access_token"`
	LastCursor        string
	UserId            string            `json:"user_id"`
}

func (authorization *Authorization) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "Authorizations",
		StringID: authorization.UserId,
	}
}

func (authorization *Authorization) Query() *datastore.Query {
	return datastore.NewQuery(authorization.KeySpec().Kind)
}

