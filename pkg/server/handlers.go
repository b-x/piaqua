package server

import (
	"net/http"
	"piaqua/pkg/controller"

	"github.com/julienschmidt/httprouter"
)

func state(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		content, err := c.GetControllerState()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}
