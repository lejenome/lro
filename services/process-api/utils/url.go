package utils

import (
	"net/http"
)

const (
	NginxForwardHostHeader = "X-Forwarded-Host"
)

/*
 ResolveHost function resolve Host from Request.
 Function code Taken from https://github.com/gin-contrib/location/blob/master/config.go
*/
func ResolveHost(r *http.Request) (host string) {
	switch {
	case r.Header.Get(NginxForwardHostHeader) != "":
		return r.Header.Get(NginxForwardHostHeader)
	case r.Header.Get("X-Host") != "":
		return r.Header.Get("X-Host")
	case r.Host != "":
		return r.Host
	case r.URL.Host != "":
		return r.URL.Host
	default:
		return ""
	}
}
