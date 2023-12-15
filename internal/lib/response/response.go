package response

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"msg,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Err(msg string) Response {
	return Response{
		Status:  StatusError,
		Message: msg,
	}
}

func SendResponse(w http.ResponseWriter, r *http.Request, status int, response interface{}) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}
