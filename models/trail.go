package models

import (
	"github.com/drborges/appx"
	"time"
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
	VideoTrail
)

func (trail *Trail) KeySpec() *appx.KeySpec {
	return &appx.KeySpec{
		Kind:       "Trails",
		StringID: trail.Revision,
	}
}
