// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package middleware
import "net/http"


type JsonResponse struct {

}

func NewJsonResponse() *JsonResponse {
	return &JsonResponse{}
}

func (l *JsonResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json")
	next(rw, r)
}
