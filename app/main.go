package app
import (
	"net/http"
	"github.com/bearchinc/trails-api"
)

func init() {
	http.Handle("/", trails.Routes())
}
