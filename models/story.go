package models
import "appengine/datastore"

type Story struct {
	Trails    []Trail `json:"trails"`
	ImagePath string `json:"image_path"`
	Title     string `json:"title"`
}

var Stories = struct {
	ByAccount func(account *Account) *datastore.Query
}{
	ByAccount: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Trail).KeySpec().Kind).
		Ancestor(account.Key()).
		Filter("Likeness=", LikedIt).
		Order("CreatedAt")
	},
}