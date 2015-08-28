package middlewares

import (
	"appengine"
	"encoding/base64"
	"github.com/bearchinc/trails-api/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

func AuthorizationAccountProvider(c appengine.Context, request *http.Request, render render.Render, martiniContext martini.Context) {
	authToken := extractAuthToken(request)

	if authToken == "" {
		render.Status(http.StatusUnauthorized)
		return
	}

	currentAccount := models.Accounts.ByAuthToken(authToken)
	if currentAccount != nil {
		render.Status(http.StatusUnauthorized)
	}

	martiniContext.Map(currentAccount)
}

func extractAuthToken(request *http.Request) string {
	if authTokenFromHeader := decodeAuthHeader(request.Header.Get("Authorization")); authTokenFromHeader != "" {
		return authTokenFromHeader
	}
	if authTokenFromQueryParams := request.FormValue("authToken"); authTokenFromQueryParams != "" {
		return authTokenFromQueryParams
	}
	return ""
}

func decodeAuthHeader(authHeaderEncoded string) string {
	if authHeaderEncoded == "" {
		return ""
	}

	authHeaderEncoded = strings.Replace(authHeaderEncoded, "Basic ", "", 1)
	auth_header, err := base64.StdEncoding.DecodeString(authHeaderEncoded)
	if err != nil {
		return ""
	}

	user_name := strings.Split(string(auth_header), ":")

	return user_name[0]
}
