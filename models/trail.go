package models

import (
	"github.com/drborges/appx"
	"time"
	"appengine/datastore"
	"math/rand"
	"appengine"
	"github.com/drborges/geocoder/providers/google"
	"appengine/urlfetch"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

type Trail struct {
	appx.Model

	Revision    string               `json:"-"`
	Path        string               `json:"media_path"`
	ThumbExists bool                 `json:"thumb_exists"`

	MimeType    string               `json:"mime_type"`
	CreatedAt   time.Time            `json:"created_at"`
	GeoPoint    appengine.GeoPoint  `json:"geo_point"`
	Tags        []string              `json:"geo_point"`
	Bytes       int64                `json:"bytes"`

	Type        TrailType            `json:"trail_type"`

	Likeness    LikenessType         `json:"likeness"`
	EvaluatedOn time.Time            `json:"evaluated_on"`
}

type LikenessType int
const (
	NotEvaluated    LikenessType = iota
	LikedIt
	DislikedIt
)

const Uncategorized = "Uncategorized"

type TrailType int
const (
	PhotoTrail    TrailType = iota
	AudioTrail
	VideoTrail
)

func (trail *Trail) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:        "Trails",
		StringID:    trail.Revision,
		HasParent:    true,
	}
}

func randomDate() time.Time {
	rand.Seed(time.Now().Unix())
	randomMonth := rand.Intn(108)
	rand.Seed(time.Now().Unix() + int64(randomMonth))
	randomDay := rand.Intn(30)

	return time.Now().AddDate(0, -randomMonth, -randomDay)
}

func likeness(trailId string, likeness LikenessType, db *appx.Datastore, context appengine.Context) error {
	trail := Trail{}
	trail.SetEncodedKey(trailId)

	if err := db.Load(&trail); err != nil {
		println("The error: ", err.Error())
		return err
	}

	trail.Likeness = likeness
	trail.EvaluatedOn = time.Now()

	if trail.Likeness == LikedIt {
		trail.Tags = fetchLatLngFromGoogle(trail, context)
		trail.storeTags(context, db)
	}

	if err := db.Save(&trail); err != nil {
		println("The error: ", err.Error())
		return err
	}

	return nil
}

func (trail *Trail) storeTags(context appengine.Context, db *appx.Datastore) {
	extractTags := func(data stream.T) stream.T {
		trail := data.(*Trail)
		return trail.AllTags()
	}

	err := rivers.FromData(trail).FlatMap(extractTags).Each(func(data stream.T) {
		tag := data.(*Tag)
		tag.SetParentKey(trail.ParentKey())
		tag.ImagePath = trail.Path
		tag.ImageProvider = DropBox
		datastore.RunInTransaction(context, func(c appengine.Context) error {
			if err := db.Load(tag); err == nil {
				tag.LikenessCount++
			} else {
				tag.LikenessCount = 1
			}

			return db.Save(tag)
		}, nil)
	}).Drain()
	if err != nil {
		println(">>>>Screwed", err.Error())
	}
}

func (trail Trail)AllTags() []*Tag {
	if len(trail.Tags) == 1 {
		return []*Tag{ trail.Unspecified() }
	} else {
		return []*Tag{
			trail.City(),
			trail.State(),
			trail.Country(),
		}
	}
}

func (trail Trail) Unspecified() *Tag {
	return &Tag{Type: TagTypeUnspecified, Value:trail.Tags[0]}
}

func (trail Trail) City() *Tag {
	return &Tag{Type: TagTypeCity, Value:trail.Tags[TagTypeCity]}
}

func (trail Trail) State() *Tag {
	return &Tag{Type: TagTypeState, Value: trail.Tags[TagTypeState]}
}

func (trail Trail) Country() *Tag {
	return &Tag{Type: TagTypeCountry, Value: trail.Tags[TagTypeCountry]}
}

func fetchLatLngFromGoogle(trail Trail, context appengine.Context) []string {
	geoCoder := &google.Geocoder{
		HttpClient:             urlfetch.Client(context),
		ReverseGeocodeEndpoint: google.ReverseGeocodeEndpoint + "&key=AIzaSyC1O6FZtjFDSJz5zCqVbVlVOr60gDYg_Zw",
	}

	if (trail.GeoPoint == appengine.GeoPoint{}) { return []string{Uncategorized} }

	println(">>>>>>About to fetch from Google!")
	res, err := geoCoder.ReverseGeocode(trail.GeoPoint.Lat, trail.GeoPoint.Lng)
	if err != nil { return []string{Uncategorized} }

	var address google.Address
	google.ReadResponse(res, &address)

	return []string{address.FullCountry, address.FullState, address.FullCity}
}

var Trails = struct {
	ByNextEvaluation func(account *Account) *datastore.Query
	ByAccount        func(account *Account) *datastore.Query
	ByTag            func(tagId string, account *Account) *datastore.Query
	Like             func(trailId string, db *appx.Datastore, context appengine.Context) error
	Dislike          func(trailId string, db *appx.Datastore, context appengine.Context) error

}{
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

	ByTag:func(tagId string, account *Account) *datastore.Query {
		return datastore.NewQuery(new(Trail).KeySpec().Kind).
		Ancestor(account.Key()).
		Filter("Tags=", tagId).
		Order("-EvaluatedOn")
	},

	Like: func(trailId string, db *appx.Datastore, context appengine.Context) error {
		return likeness(trailId, LikedIt, db, context)
	},

	Dislike: func(trailId string, db *appx.Datastore, context appengine.Context) error {
		return likeness(trailId, DislikedIt, db, context)
	},
}
