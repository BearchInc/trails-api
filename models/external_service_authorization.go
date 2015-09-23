package models

import (
	"github.com/drborges/appx"
	"strconv"
	"appengine/datastore"
)

type AuthorizationType int

const Url AuthorizationType = 0
const DropBox AuthorizationType = 1


func (t AuthorizationType) String() string {
	return strconv.Itoa(int(t))
}

type ExternalServiceAuthorization struct {
	appx.Model

	AuthorizationType AuthorizationType
	AccessToken       string
	LastCursor        string
	UserId            string
}

func (authorization *ExternalServiceAuthorization) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "ExternalServiceAuthorization",
		StringID:   authorization.UserId,
		HasParent:  true,
	}
}

func (authorization *ExternalServiceAuthorization) Query() *datastore.Query {
	return datastore.NewQuery(authorization.KeySpec().Kind).Ancestor(authorization.ParentKey())
}

