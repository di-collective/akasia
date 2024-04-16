package common

import "net/http"

type REST interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}
