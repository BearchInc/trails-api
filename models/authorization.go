package models

import (
	"github.com/drborges/appx"
	"strconv"
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
	UserId				string			`json:"user_id"`
}

func (account *Authorization) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "Authorizations",
		Incomplete: true,
		StringID: account.UserId,
	}
}
