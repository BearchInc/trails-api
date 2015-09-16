package models

import (
	"github.com/drborges/appx"
	"time"
	"appengine/datastore"
"math/rand"
)

type Trail struct {
	appx.Model

	Revision    string	`json:"-"`
	Path        string	`json:"media_path"`
	ThumbExists bool	`json:"thumb_exists"`

	MimeType    string    `json:"mime_type"`
	CreatedAt   time.Time `json:"created_at"`
	Location    []float64 `json:"location"`
	Bytes       int64        `json:"bytes"`

	Type        TrailType	`json:"trail_type"`

	Likeness    LikenessType	`json:"likeness"`
}

type LikenessType int
const (
	NotEvaluated    LikenessType = iota
	LikedIt
	DislikedIt
)

type TrailType int
const (
	PhotoTrail    TrailType = iota
	AudioTrail
	VideoTrail
)

func (trail *Trail) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "Trails",
		StringID: trail.Revision,
	}
}

func randomDate() time.Time {
	rand.Seed(time.Now().Unix())
	randomMonth := rand.Intn(18)
	rand.Seed(time.Now().Unix() + int64(randomMonth))
	randomDay := rand.Intn(30)

	return time.Now().AddDate(0, -randomMonth, -randomDay)
}

func likeness(trailId string, likeness LikenessType, db *appx.Datastore) error {
	trail := Trail{}
	trail.SetEncodedKey(trailId)

	if err := db.Load(&trail); err != nil {
		println("The error: ", err.Error())
		return err
	}

	trail.Likeness = likeness

	if err := db.Save(&trail); err != nil {
		println("The error: ", err.Error())
		return err
	}

	return nil
}

var Trails = struct {
	ByNextEvaluation 	func(account *Account) *datastore.Query
	ByAccount		 	func(account * Account) *datastore.Query
	Like 				func(trailId string, db *appx.Datastore) error
	Dislike 			func(trailId string, db *appx.Datastore) error

} {
	ByNextEvaluation: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Trail).KeySpec().Kind).
			Ancestor(account.Key()).
			Filter("CreatedAt >", randomDate()).
			Filter("Likeness =", NotEvaluated).
			Order("CreatedAt").
			Limit(6)
	},

	ByAccount: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Trail).KeySpec().Kind).
		Ancestor(account.Key())
	},

	Like: func(trailId string, db *appx.Datastore) error {
		return likeness(trailId, LikedIt, db)
	},

	Dislike: func(trailId string, db *appx.Datastore) error {
		return likeness(trailId, DislikedIt, db)
	},
}