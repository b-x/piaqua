package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"piaqua/pkg/controller"
	"piaqua/pkg/model"
	"strconv"

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
		_, _ = w.Write(content)
	}
}

func addAction(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		var action model.Action
		err := decoder.Decode(&action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := c.AddAction(&action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ret, _ := json.Marshal(struct {
			Id int `json:"id"`
		}{id})
		_, _ = w.Write(ret)
	}
}

func updateAction(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var action model.Action
		err = decoder.Decode(&action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.UpdateAction(id, &action)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func toggleAction(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.ToggleAction(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func removeAction(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.RemoveAction(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func setSensorName(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.SetSensorName(id, string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func setRelayName(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.SetRelayName(id, string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func addRelayTask(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		relayId, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var task model.RelayTask
		err = decoder.Decode(&task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := c.AddRelayTask(relayId, &task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ret, _ := json.Marshal(struct {
			Id int `json:"id"`
		}{id})
		_, _ = w.Write(ret)
	}
}

func updateRelayTask(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		relayId, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		taskId, err := strconv.Atoi(p.ByName("tid"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var task model.RelayTask
		err = decoder.Decode(&task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.UpdateRelayTask(relayId, taskId, &task)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func removeRelayTask(c *controller.Controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		relayId, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		taskId, err := strconv.Atoi(p.ByName("tid"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = c.RemoveRelayTask(relayId, taskId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
