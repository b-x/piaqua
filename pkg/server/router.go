package server

import (
	"net/http"
	"piaqua/pkg/config"
	"piaqua/pkg/controller"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var apiRoutes = []struct {
	method  string
	path    string
	handler func(*controller.Controller) httprouter.Handle
}{
	{"GET", "/api/state", state},
	{"POST", "/api/action", addAction},
	{"PUT", "/api/action/:id", updateAction},
	{"PUT", "/api/action/:id/toggle", toggleAction},
	{"DELETE", "/api/action/:id", removeAction},
	{"PUT", "/api/sensor/:id/name", setSensorName},
	{"PUT", "/api/relay/:id/name", setRelayName},
	{"POST", "/api/relay/:id/task", addRelayTask},
	{"PUT", "/api/relay/:id/task/:tid", updateRelayTask},
	{"DELETE", "/api/relay/:id/task/:tid", removeRelayTask},
}

func newHandler(conf *config.ServerConf, c *controller.Controller) http.Handler {

	apiRouter := &httprouter.Router{}
	for _, r := range apiRoutes {
		apiRouter.Handle(r.method, r.path, basicAuth(conf.Credentials, r.handler(c)))
	}

	return &muxHandler{
		api:    apiRouter,
		static: http.FileServer(http.Dir(conf.Path)),
	}
}

type muxHandler struct {
	http.Handler
	api    http.Handler
	static http.Handler
}

func (h *muxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		h.api.ServeHTTP(w, r)
		return
	}
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, nosnippet, noarchive")
	h.static.ServeHTTP(w, r)
}

func basicAuth(credentials map[string]string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()

		if hasAuth {
			pass, found := credentials[user]
			if found && password == pass {
				h(w, r, ps)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
