package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/franela/play-with-docker/config"
	"github.com/franela/play-with-docker/services"
)

type NewSessionResponse struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
}

func NewSession(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if !services.IsHuman(req) {
		// User it not a human
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Only humans are allowed!"))
		return
	}

	reqDur := req.Form.Get("session-duration")

	duration := services.GetDuration(reqDur)
	s, err := services.NewSession(duration)
	if err != nil {
		log.Println(err)
		//TODO: Return some error code
	} else {

		hostname := fmt.Sprintf("%s.%s", config.PWDCName, req.Host)
		// If request is not a form, return sessionId in the body
		if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			resp := NewSessionResponse{SessionId: s.Id, Hostname: hostname}
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(resp)
			return
		}
		http.Redirect(rw, req, fmt.Sprintf("http://%s/p/%s", hostname, s.Id), http.StatusFound)
	}
}
