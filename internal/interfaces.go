package internal

import "net/http"

type Service interface {
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
