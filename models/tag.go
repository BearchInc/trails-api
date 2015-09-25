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

const TagTypeUnspecified TagType = TagTypeCity

type Tag struct {
	appx.Model
	Type          TagType   `json:"-"`
	Value         string    `json:"title"`
	LikenessCount int        `json:"likeness_count"`
	ImagePath     string    `json:"image_path"`
	ImageProvider AuthorizationType `json:"authorization_type"`
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
}{
	ByAccount: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Tag).KeySpec().Kind).
		Ancestor(account.Key()).
		Filter("Type=", TagTypeCity).
		Order("-LikenessCount")
	},
}