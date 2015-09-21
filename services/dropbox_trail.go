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
	"appengine"
)

func DropboxDelta(req *http.Request,ds *appx.Datastore, authorization *models.ExternalServiceAuthorization, existingItem stream.PredicateFn) {
	rivers.DebugEnabled = true

	dropboxClient := dropboxClient(newappengine.NewContext(req), authorization.AccessToken)
	builder := DropboxDeltaProducerBuilder{Client: dropboxClient, CurrentCursor: authorization.LastCursor}

	deltaStream, debugStream := rivers.
				From(builder.Build()).
				Drop(nonMediaFiles).
				Map(toTrailForAuthorization(authorization)).
				Drop(existingItem).
				Split()

	println("About to split data")

	err := deltaStream.BatchBy(&appx.DatastoreBatch{Size: 500}).
				Each(saveBatch(ds)).
				Drain()
	println(fmt.Sprintf("Error while batch saving the data: %v", err))

	total, errTotal := debugStream.Reduce(0, func(sum, item stream.T) stream.T{
		return sum.(int) + 1
	}).CollectFirst()

	if errTotal != nil {
		println(fmt.Sprintf("########## Guaging error %v", err))
	}
	println(fmt.Sprintf("Total entries: %v", total.(int)))

	authorization.LastCursor = builder.CurrentCursor
	ds.Save(authorization)

	if err != nil {
		println(fmt.Sprintf("########## Delta error %v", err))
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
	return true
}

func toTrailForAuthorization(authorization *models.ExternalServiceAuthorization) stream.MapFn {
	return func(data stream.T) stream.T {
		trailItem := toTrail(data)
		trail := trailItem.(*models.Trail)
		trail.SetParentKey(authorization.ParentKey())

		return trail
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
		GeoPoint: geoPointFrom(item),
	}

	return trail
}

func geoPointFrom(item *dropbox.Entry) *appengine.GeoPoint {
	if item.PhotoInfo != nil {
		return &appengine.GeoPoint{Lat: item.PhotoInfo.LatLong[0], Lng: item.PhotoInfo.LatLong[1]}
	}

	if item.VideoInfo != nil {
		return &appengine.GeoPoint{Lat: item.VideoInfo.LatLong[0], Lng: item.VideoInfo.LatLong[1]}
	}

	return nil
}

func saveBatch(ds *appx.Datastore) stream.EachFn {
	return func(data stream.T) {
		trails := data.(*appx.DatastoreBatch)

		if err := ds.Save(trails.Items...); err != nil {
			panic(err)
		}
	}
}
