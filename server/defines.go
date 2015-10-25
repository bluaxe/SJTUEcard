package server

import (
	"net/http"
	"net/url"
)

type handler func(url.Values, http.ResponseWriter)
