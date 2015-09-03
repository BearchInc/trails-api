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

func DropboxDelta(req *http.Request, ds *appx.Datastore, authorization *models.Authorization) {
	steamContext := rivers.NewContext()
	dropboxClient := dropboxClient(newappengine.NewContext(req), authorization.AccessToken)
	builder := DropboxDeltaProducerBuilder{Context: steamContext, Client: dropboxClient, CurrentCursor: authorization.LastCursor}

	err := rivers.NewWith(builder.Context).
				From(builder.Build()).
				Drop(notMedia).
				Map(toTrail).
				Drop(alreadyCategorized(ds)).
				BatchBy(&appx.DatastoreBatch{Size: 500}).
				Each(saveBatch(ds)).
				Drain()

	authorization.LastCursor = builder.CurrentCursor
	ds.Save(authorization)

	println(fmt.Sprintf("########## %v", err))
}

func DropboxInit(req *http.Request, ds *appx.Datastore, authorization *models.Authorization) {

//	accessToken := "4GFHgi3IL2IAAAAAAAAAClnkF5SyJKu7pcC64oNVg18r6tecpNkJdyY7spxDw9hF"
//	accessToken := "Cdn2AKzQYPwAAAAAAAAT-9gnnBDOG68LO1Oms5qW9-4WLme9TX02EHLiQLgR2qoA"

	steamContext := rivers.NewContext()
	dropboxClient := dropboxClient(newappengine.NewContext(req), authorization.AccessToken)
	builder := DropboxDeltaProducerBuilder{Context: steamContext, Client: dropboxClient}

	err := rivers.NewWith(builder.Context).
					From(builder.Build()).
					Drop(notMedia).
					Map(toTrail).
					BatchBy(&appx.DatastoreBatch{Size: 500}).
					Each(saveBatch(ds)).
					Drain()

	authorization.LastCursor = builder.CurrentCursor
	ds.Save(authorization)

	println(fmt.Sprintf("########## %v", builder.CurrentCursor))
	println(fmt.Sprintf("########## %v", err))
}

func notMedia(data stream.T) bool {
	item := data.(*dropbox.Entry)
	if strings.HasPrefix(item.MimeType, "image") { return false }
	if strings.HasPrefix(item.MimeType, "video") { return false }

	return true
}

func alreadyCategorized(db *appx.Datastore) stream.PredicateFn {

	return func(data stream.T) bool {
		trail := data.(*models.Trail)

		err := db.Load(trail)

		return err == nil
	}
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