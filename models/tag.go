package models
import (
"github.com/drborges/appx"
"appengine/datastore"
)

type TagType int
const (
	TagTypeCountry    TagType = iota
	TagTypeState
	TagTypeCity
)

type Tag struct {
	appx.Model
	Type          TagType
	Value         string
	LikenessCount int
}

func (tag *Tag) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:        "Tags",
		StringID:    tag.Value,
		HasParent:  true,
	}
}

var Tags = struct {
	ByAccount func(account *Account) *datastore.Query
} {
	ByAccount: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Tag).KeySpec().Kind).
			Ancestor(account.Key()).
			Filter("Type=", TagTypeCity).
			Order("-LikenessCount")
	},
}