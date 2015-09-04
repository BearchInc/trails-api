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

func DropboxDelta(req *http.Request, ds *appx.Datastore, authorization *models.ExternalServiceAuthorization, existingItem stream.PredicateFn) {
	rivers.DebugEnabled = true

	steamContext := rivers.NewContext()
	dropboxClient := dropboxClient(newappengine.NewContext(req), authorization.AccessToken)
	builder := DropboxDeltaProducerBuilder{Context: steamContext, Client: dropboxClient, CurrentCursor: authorization.LastCursor}

	deltaStream, debugStream := rivers.NewWith(builder.Context).
				From(builder.Build()).
				Drop(nonMediaFiles).
				Map(toTrail).
				Drop(existingItem).
				Split()


	err := deltaStream.BatchBy(&appx.DatastoreBatch{Size: 500}).
				Each(saveBatch(ds)).
				Drain()

	total, _ := debugStream.Reduce(0, func(sum, item stream.T) stream.T{
		return sum.(int) + 1
	}).CollectFirst()

	println(fmt.Sprintf("Total entries: %v", total.(int)))

	authorization.LastCursor = builder.CurrentCursor
	ds.Save(authorization)

	if err != nil {
		println(fmt.Sprintf("########## %v", err))
	}
}

func nonMediaFiles(data stream.T) bool {
	item := data.(*dropbox.Entry)

	println(fmt.Sprintf("Mime: %v AND Name: %v", item.MimeType, item.Path))

	if strings.HasPrefix(item.MimeType, "image") {
		println("")
		return false
	}
	if strings.HasPrefix(item.MimeType, "video") {
		println("")
		return false
	}

	println("Skipped the file")
	println("")
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