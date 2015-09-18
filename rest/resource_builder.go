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

func FromTrails(trails []*models.Trail) []stream.T {
	rivers.DebugEnabled = true

	if len(trails) == 0 {
		return []stream.T{}
	}

	resources, err := rivers.FromSlice(trails).
						Map(toTrailResource).
						Collect()
	if err != nil {
		panic(err)
	}
	return resources
}

func toTrailResource(item stream.T) stream.T {
	trail := item.(*models.Trail)
	selfPath := fmt.Sprint("/account/trails/", trail.EncodedKey())

	return TrailResource{
		Trail: trail,
		LikePath: selfPath + "/like",
		DislikePath: selfPath + "/dislike",
	}
}
