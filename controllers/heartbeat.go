package controllers

import (
	"net/http"
)

type HeartbeatController struct {
	ResponseController
}

func (c *HeartbeatController) Get(w http.ResponseWriter, r *http.Request) {
	c.respondWithJSON(w, http.StatusOK, "❤")
}
