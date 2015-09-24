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

type TagResource struct {
	Title         string `json:"title"`
	LikenessCount int `json:"likeness_count"`
	ImagePath		string `json:"image_path"`
	AuthorizationType models.AuthorizationType `json:"authorization_type"`

	SelfPath      string `json:"self_path"`
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
	selfPath := fmt.Sprintf("/account/trails/%v", trail.EncodedKey())

	return TrailResource{
		Trail: trail,
		LikePath: selfPath + "/like",
		DislikePath: selfPath + "/dislike",
	}
}

func FromTags(tags []*models.Tag) []stream.T {
	rivers.DebugEnabled = true

	if len(tags) == 0 {
		return []stream.T{}
	}

	resources, err := rivers.FromSlice(tags).
	Map(toTagResource).
	Collect()

	if err != nil {
		panic(err)
	}
	return resources
}

func toTagResource(item stream.T) stream.T {
	tag := item.(*models.Tag)
	selfPath := fmt.Sprintf("/account/trails/%v", tag.Value)

	return TagResource{
		Title: tag.Value,
		LikenessCount: tag.LikenessCount,
		ImagePath: "http://allworldtowns.com/data_images/countries/hawaii/hawaii-09.html",
		AuthorizationType: models.Url,
		SelfPath: selfPath,
	}
}