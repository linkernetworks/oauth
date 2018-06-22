package endpoint

import (
	"net/http"
)

func (ep *Endpoint) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
