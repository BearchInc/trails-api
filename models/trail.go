package models
import (
	"github.com/drborges/appx"
	"time"
)

type Trail struct {
	appx.Model

	Path 		string			`json:"path"`
	ThumbExists bool			`json:"thumb_exists"`

	MimeType	string			`json:"mime_type"`
	CreatedAt 	time.Time		`json:"created_at"`
	Location 	time.Location	`json:"location"`
	Bytes		int				`json:"bytes"`
}
