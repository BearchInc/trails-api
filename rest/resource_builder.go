package rest
import (
	"github.com/bearchinc/trails-api/models"
	"fmt"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

type TrailResource struct {
	*models.Trail

	LikePath    string `json:"like_path"`
	DislikePath string `json:"dislike_path"`
}

type PagedResources struct {
	Resources []stream.T	`json:"resources"`
	NextPage string    `json:"next_page"`
}

func FromTrails(trails []*models.Trail, cursor string) PagedResources {
	var resources PagedResources
	rivers.DebugEnabled = true
	rivers.FromSlice(trails).
		Map(toTrailResource).
		Batch(len(trails)).
		Map(toTrailsResources(cursor)).
		CollectFirstAs(&resources)
	return resources
}

func toTrailResource(item stream.T) stream.T {
	trail := item.(*models.Trail)
	selfPath := fmt.Sprint("/accounts/trails/%v", trail.EncodedKey())

	return TrailResource{
		Trail: trail,
		LikePath: selfPath + "/like",
		DislikePath: selfPath + "/dislike",
	}
}

func toTrailsResources(cursor string) stream.MapFn {
	return func(item stream.T) stream.T {
		trailsResources := item.([]stream.T)
		return PagedResources{
			Resources: trailsResources,
			NextPage: cursor,
		}
	}

}