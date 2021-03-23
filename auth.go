package shu

import "net/http"

// Auth auth
type Auth interface {
	Auth(request *http.Request) error
}

// BasicAuth basic auth
type BasicAuth struct {
	// Name name
	Name string
	// Password password
	Password string
}

func (auth *BasicAuth) Auth(r http.Request) error {
	r.SetBasicAuth(auth.Name, auth.Password)
	return nil
}
