package rest
import (
	"github.com/bearchinc/trails-api/models"
	"fmt"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"net/url"
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
	return from(trails, toTrailResource)
}

func FromTags(tags []*models.Tag) []stream.T {
	return from(tags, toTagResource)
}

func from(from stream.T, mapper func(item stream.T) stream.T) []stream.T {
	resources, _ := rivers.FromSlice(from).
		Map(mapper).
		Collect()

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


func toTagResource(item stream.T) stream.T {
	tag := item.(*models.Tag)
	selfPath := fmt.Sprintf("/account/trails/tags/%v", tag.Value)
	url, _ := url.Parse(selfPath)
	selfPath = url.RequestURI()

	return TagResource{
		Title: tag.Value,
		LikenessCount: tag.LikenessCount,
		ImagePath: "http://allworldtowns.com/data_images/countries/hawaii/hawaii-09.jpg",
		AuthorizationType: models.Url,
		SelfPath: selfPath,
	}
}