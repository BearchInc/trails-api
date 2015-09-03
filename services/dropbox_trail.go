package services
import (
	"github.com/drborges/rivers"
	newappengine "google.golang.org/appengine"
	"github.com/drborges/rivers/stream"
	"fmt"
	"net/http"
	"github.com/stacktic/dropbox"
	"strings"
	"github.com/bearchinc/trails-api/models"
	"time"
	"github.com/drborges/appx"
)

func DropboxDeltaFirstTime(req *http.Request, ds *appx.Datastore, authorization *models.Authorization) {

//	accessToken := "4GFHgi3IL2IAAAAAAAAAClnkF5SyJKu7pcC64oNVg18r6tecpNkJdyY7spxDw9hF"
//	accessToken := "Cdn2AKzQYPwAAAAAAAAT-9gnnBDOG68LO1Oms5qW9-4WLme9TX02EHLiQLgR2qoA"

	steamContext := rivers.NewContext()
	db := DropboxBuilder{Context: steamContext}

	dropboxClient := dropboxClient(newappengine.NewContext(req), authorization.AccessToken)

	producer := db.DropboxProducer(dropboxClient)

	err := rivers.NewWith(steamContext).
					From(producer).
					DropIf(notMedia).
					Map(toTrail).
					BatchBy(&appx.DatastoreBatch{Size: 500}).
					Each(saveBatch(ds)).
					Drain()

	authorization.LastCursor = db.CurrentCursor
	ds.Save(authorization)

	println(fmt.Sprintf("########## %v", db.CurrentCursor))
	println(fmt.Sprintf("########## %v", err))
}

func notMedia(data stream.T) bool {
	item := data.(*dropbox.Entry)
	if strings.HasPrefix(item.MimeType, "image") { return false }
	if strings.HasPrefix(item.MimeType, "video") { return false }

	return true
}

func toTrail(data stream.T) stream.T {
	item := data.(*dropbox.Entry)

	trail := &models.Trail{
		Path: item.Path,
		ThumbExists: item.ThumbExists,
		MimeType: item.MimeType,
		CreatedAt: time.Time(item.ClientMtime),
		Bytes: item.Bytes,
		Likeness: models.NotEvaluated,
		Revision: item.Revision,
	}

	if item.PhotoInfo != nil { trail.Location = item.PhotoInfo.LatLong }
	if item.VideoInfo != nil { trail.Location = item.VideoInfo.LatLong }

	return trail
}

func saveBatch(ds *appx.Datastore) stream.EachFn {
	return func(data stream.T) {
		trails := data.(*appx.DatastoreBatch)

		if err := ds.Save(trails.Items...); err != nil {
			panic(err)
		}
	}
}