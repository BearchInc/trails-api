package models

import (
	"github.com/drborges/appx"
	"time"
	"appengine/datastore"
"math/rand"
)

type Trail struct {
	appx.Model

	Revision    string
	Path        string `json:"path"`
	ThumbExists bool   `json:"thumb_exists"`

	MimeType    string    `json:"mime_type"`
	CreatedAt   time.Time `json:"created_at"`
	Location    []float64 `json:"location"`
	Bytes       int64        `json:"bytes"`

	Type        TrailType

	Likeness    LikenessType
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
	randomMonth := rand.Intn(36)

	return time.Now().AddDate(0, -randomMonth, 0)
}

var Trails = struct {
	ByNextEvaluation func(account *Account) *datastore.Query
} {
	ByNextEvaluation: func(account *Account) *datastore.Query {
		return datastore.NewQuery(new(Trail).KeySpec().Kind).
			Ancestor(account.Key()).
			Filter("CreatedAt >", randomDate()).
			Order("CreatedAt").
			Limit(6)
	},
}