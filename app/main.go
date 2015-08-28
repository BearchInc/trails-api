package app

import (
	"github.com/bearchinc/trails-api"
	"net/http"
)

func init() {
	http.Handle("/", trails.Routes())
}
